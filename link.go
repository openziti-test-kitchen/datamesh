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
	Destination
	Start() error
	Peer() *identity.TokenId
	Direction() LinkDirection
	SendControl(c *Control) error
}

type link struct {
	cfg           *LinkConfig
	addr          Address
	id            *identity.TokenId
	peer          *identity.TokenId
	ch            channel.Channel
	direction     LinkDirection
	pingResponses chan *channel.Message
}

func newLink(cfg *LinkConfig, direction LinkDirection) *link {
	l := &link{
		cfg:           cfg,
		direction:     direction,
		pingResponses: make(chan *channel.Message, cfg.PingQueueLength),
	}
	return l
}

func (self *link) setChannel(ch channel.Channel) {
	self.ch = ch
	self.id = &identity.TokenId{Token: ch.ConnectionId()}
}

func (self *link) Address() Address {
	return Address(self.id.Token)
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

func (self *link) Start() error {
	go self.pinger()
	return nil
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

type linkBindHandler struct {
	datamesh *Datamesh
	link     *link
}

func newLinkBindHandler(datamesh *Datamesh, link *link) *linkBindHandler {
	return &linkBindHandler{datamesh, link}
}

func (self *linkBindHandler) BindChannel(ch channel.Channel) error {
	ch.AddReceiveHandler(&linkControlReceiveHandler{self.link})
	ch.AddReceiveHandler(&linkPayloadReceiveHandler{self.datamesh, self.link})
	ch.AddReceiveHandler(&linkAcknowledgementReceiveHandler{self.datamesh, self.link})
	return nil
}

type linkControlReceiveHandler struct {
	link *link
}

func (_ *linkControlReceiveHandler) ContentType() int32 {
	return int32(ControlContentType)
}

func (self *linkControlReceiveHandler) HandleReceive(msg *channel.Message, ch channel.Channel) {
	log := pfxlog.ContextLogger(ch.ConnectionId())
	if ctrl, err := UnmarshalControl(msg); err == nil {
		if ctrl.Flags == uint32(PingRequestControlFlag) {
			var found bool
			var pingId string
			pingId, found = msg.GetStringHeader(PingIdHeaderKey)
			if found {
				var stamp uint64
				stamp, found = msg.GetUint64Header(PingTimestampHeaderKey)
				if found {
					headers := newHeaders()
					headers.PutBytes(PingIdHeaderKey, []byte(pingId))
					headers.PutInt64(PingTimestampHeaderKey, int64(stamp))
					ctrl := NewControl(uint32(PingResponseControlFlag), headers)
					err := self.link.SendControl(ctrl)
					if err == nil {
						logrus.Infof("sent response to [ping/%s]", pingId)
					} else {
						logrus.Errorf("error responding [ping/%s] (%v)", pingId, err)
					}
				} else {
					logrus.Errorf("missing timestamp")
				}
			} else {
				logrus.Errorf("missing ping identity")
			}
		} else if ctrl.Flags == uint32(PingResponseControlFlag) {
			self.link.pingResponses <- msg
		} else {
			log.Error("unknown flags")
		}
	} else {
		log.Errorf("error unmarshaling control message (%v)", err)
	}
}

type linkPayloadReceiveHandler struct {
	datamesh *Datamesh
	l        *link
}

func (_ *linkPayloadReceiveHandler) ContentType() int32 {
	return int32(PayloadContentType)
}

func (self *linkPayloadReceiveHandler) HandleReceive(msg *channel.Message, ch channel.Channel) {
	log := pfxlog.ContextLogger(ch.ConnectionId())
	if p, err := UnmarshalPayload(msg); err == nil {
		if err := self.datamesh.Fwd.ForwardPayload(p.CircuitId, self.l.Address(), p); err != nil {
			log.Errorf("error forwarding (%v)", err)
		}
	} else {
		log.Errorf("error unmarshalling (%v)", err)
	}
}

type linkAcknowledgementReceiveHandler struct {
	datamesh *Datamesh
	l        *link
}

func (_ *linkAcknowledgementReceiveHandler) ContentType() int32 {
	return int32(AckContentType)
}

func (_ *linkAcknowledgementReceiveHandler) HandleReceive(m *channel.Message, _ channel.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}
