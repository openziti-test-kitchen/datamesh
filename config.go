package datamesh

import "github.com/openziti/foundation/transport"

type Config struct {
	LinkListeners []*LinkListenerConfig
	LinkDialers   []*LinkDialerConfig
}

type LinkListenerConfig struct {
	Id            string
	BindAddress   transport.Address
	Advertisement transport.Address
}

type LinkDialerConfig struct {
	Id          string
	BindAddress transport.Address
}
