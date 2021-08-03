package datamesh

import (
	"github.com/openziti-incubator/datamesh/channel"
	cmap "github.com/orcaman/concurrent-map"
)

type Forwarder struct {
	table cmap.ConcurrentMap
}

func newForwarder() *Forwarder {
	return &Forwarder{
		table: cmap.New(),
	}
}

func (fw *Forwarder) Forward(msg *channel.Message) error {
	return nil
}