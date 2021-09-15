package datamesh

type Address string

type Destination interface {
	Address() Address
	FromNetwork(data *Data) error
	Close() error
}
