package datamesh

type Address string

type Destination interface {
	Address() Address
	SendPayload(p *Payload) error
	SendAcknowledgement(a *Acknowledgement) error
	Close() error
}
