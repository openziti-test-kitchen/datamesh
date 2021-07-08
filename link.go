package datamesh

import "github.com/openziti/foundation/identity/identity"

type Link interface {
	Id() *identity.TokenId
	Peer() *identity.TokenId
	SendPayload(p *Payload) error
	SendAcknowledgement(a *Acknowledgement) error
	Close() error
}
