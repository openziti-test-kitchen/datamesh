package datamesh

import (
	"github.com/openziti/foundation/transport"
	"time"
)

type Config struct {
	Listeners []*ListenerConfig
	Dialers   []*DialerConfig
	Profile   interface{}
}

type LinkConfig struct {
	PingPeriod      time.Duration
	PingQueueLength int
	MTU             uint32
}

func LinkConfigDefaults() *LinkConfig {
	return &LinkConfig{
		PingPeriod:      time.Duration(2) * time.Second,
		PingQueueLength: 128,
		MTU:             64 * 1024,
	}
}

type ListenerConfig struct {
	Id            string
	BindAddress   transport.Address
	Advertisement transport.Address
	LinkConfig    *LinkConfig
}

func ListenerConfigDefaults() *ListenerConfig {
	return &ListenerConfig{
		LinkConfig: LinkConfigDefaults(),
	}
}

type DialerConfig struct {
	Id          string
	BindAddress transport.Address
	LinkConfig  *LinkConfig
}

func DialerConfigDefaults() *DialerConfig {
	return &DialerConfig{
		LinkConfig: LinkConfigDefaults(),
	}
}
