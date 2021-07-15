package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
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
	ch        channel2.Channel
	id        *identity.TokenId
	peer      *identity.TokenId
	direction LinkDirection
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

func (self *link) pinger() {
	log := pfxlog.ContextLogger(self.Id().Token)
	log.Info("started")
	defer log.Warn("exited")

	for {
		time.Sleep(5 * time.Second)

		if err := self.SendControl(NewControl(uint32(PingRequestControlFlag), nil)); err == nil {
			log.Info("sent control ping request")
		} else {
			log.Errorf("error sending control ping request (%v)", err)
		}
	}
}

type linkBindHandler struct{}

func newLinkBindHandler(l *link) *linkBindHandler {
	return &linkBindHandler{}
}

func (_ *linkBindHandler) BindChannel(ch channel2.Channel) error {
	ch.AddReceiveHandler(&linkControlReceiveHandler{})
	ch.AddReceiveHandler(&linkPayloadReceiveHandler{})
	ch.AddReceiveHandler(&linkAcknowledgementReceiveHandler{})
	return nil
}

type linkControlReceiveHandler struct{}

func (_ *linkControlReceiveHandler) ContentType() int32 {
	return int32(ControlContentType)
}

func (_ *linkControlReceiveHandler) HandleReceive(msg *channel2.Message, ch channel2.Channel) {
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

func (_ *linkPayloadReceiveHandler) HandleReceive(m *channel2.Message, _ channel2.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}

type linkAcknowledgementReceiveHandler struct{}

func (_ *linkAcknowledgementReceiveHandler) ContentType() int32 {
	return int32(AckContentType)
}

func (_ *linkAcknowledgementReceiveHandler) HandleReceive(m *channel2.Message, _ channel2.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}
