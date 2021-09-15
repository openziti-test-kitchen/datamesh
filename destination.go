package datamesh

type Address string

type Destination interface {
	Address() Address
	FromNetwork(data []byte) error
	Close() error
}
