package datamesh

type Address string

type Destination interface {
	Address() Address
	SendData(data *Data) error
	Close() error
}
