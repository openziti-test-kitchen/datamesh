package datamesh

type Address string

type Destination interface {
	Address() Address
	FromNetwork(data *Payload) error
	Close() error
}
