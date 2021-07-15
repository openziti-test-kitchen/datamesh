package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/pkg/errors"
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
