package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti-incubator/datamesh/channel"
	"github.com/openziti/foundation/identity/identity"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LinkDirection int8

const (
	InboundLink  LinkDirection = 0
	OutboundLink LinkDirection = 1
)

type Link interface {
	Id() *identity.TokenId
	Peer() *identity.TokenId
	Direction() LinkDirection
	SendControl(c *Control) error
	SendPayload(p *Payload) error
	SendAcknowledgement(a *Acknowledgement) error
	Close() error
}

type link struct {
	cfg           *LinkConfig
	id            *identity.TokenId
	peer          *identity.TokenId
	ch            channel.Channel
	direction     LinkDirection
	pingResponses chan *channel.Message
}

func newLink(cfg *LinkConfig, id, peer *identity.TokenId, ch channel.Channel, direction LinkDirection) *link {
	l := &link{
		cfg:           cfg,
		id:            id,
		peer:          peer,
		ch:            ch,
		direction:     direction,
		pingResponses: make(chan *channel.Message, cfg.PingQueueLength),
	}
	go l.pinger()
	return l
}

func (self *link) Id() *identity.TokenId {
	return self.id
}

func (self *link) Peer() *identity.TokenId {
	return self.peer
}

func (self *link) Direction() LinkDirection {
	return self.direction
}

func (self *link) SendControl(c *Control) error {
	return self.ch.Send(c.Marshal())
}

func (self *link) SendPayload(p *Payload) error {
	return errors.New("not implemented")
}

func (self *link) SendAcknowledgement(a *Acknowledgement) error {
	return errors.New("not implemented")
}

func (self *link) Close() error {
	return errors.New("not implemented")
}

type linkBindHandler struct{}

func (_ *linkBindHandler) BindChannel(ch channel.Channel) error {
	ch.AddReceiveHandler(&linkControlReceiveHandler{})
	ch.AddReceiveHandler(&linkPayloadReceiveHandler{})
	ch.AddReceiveHandler(&linkAcknowledgementReceiveHandler{})
	return nil
}

type linkControlReceiveHandler struct{}

func (_ *linkControlReceiveHandler) ContentType() int32 {
	return int32(ControlContentType)
}

func (_ *linkControlReceiveHandler) HandleReceive(msg *channel.Message, ch channel.Channel) {
	log := pfxlog.ContextLogger(ch.ConnectionId())
	if ctrl, err := UnmarshallControl(msg); err == nil {
		if ctrl.Flags == uint32(PingRequestControlFlag) {
			log.Info("received ping request")
		} else {
			log.Error("unknown flags")
		}
	} else {
		log.Errorf("error unmarshaling control message (%v)", err)
	}
}

type linkPayloadReceiveHandler struct{}

func (_ *linkPayloadReceiveHandler) ContentType() int32 {
	return int32(PayloadContentType)
}

func (_ *linkPayloadReceiveHandler) HandleReceive(m *channel.Message, _ channel.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}

type linkAcknowledgementReceiveHandler struct{}

func (_ *linkAcknowledgementReceiveHandler) ContentType() int32 {
	return int32(AckContentType)
}

func (_ *linkAcknowledgementReceiveHandler) HandleReceive(m *channel.Message, _ channel.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}
