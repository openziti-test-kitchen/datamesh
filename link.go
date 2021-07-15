package datamesh

import (
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/pkg/errors"
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

func (self *link) SendPayload(p *Payload) error {
	return errors.New("not implemented")
}

func (self *link) SendAcknowledgement(a *Acknowledgement) error {
	return errors.New("not implemented")
}

func (self *link) Close() error {
	return errors.New("not implemented")
}
