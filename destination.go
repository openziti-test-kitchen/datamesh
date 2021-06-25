package datamesh

import "github.com/openziti/foundation/identity/identity"

type Destination interface {
	Id() *identity.TokenId
	SendPayload(p *Payload) error
	SendAcknowledgement(a *Acknowledgement) error
	Close() error
}
