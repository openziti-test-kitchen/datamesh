package datamesh

type Endpoint interface {
	Rx([]byte) error
}

