package datamesh

import (
	"github.com/openziti-incubator/datamesh/channel"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/pkg/errors"
)

type Forwarder struct {
	table cmap.ConcurrentMap // [Address]Destination
}

func newForwarder() *Forwarder {
	return &Forwarder{
		table: cmap.New(),
	}
}

func (fw *Forwarder) addDestination(d Destination) {
	fw.table.Set(string(d.Address()), d)
}

func (fw *Forwarder) forward(srcAddr Address, msg *channel.Message) error {
	switch msg.ContentType {
	default:
		return errors.Errorf("cannot forward content type [%d]", msg.ContentType)
	}
}
