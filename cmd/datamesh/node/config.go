package node

import (
	"github.com/openziti-incubator/datamesh"
	"github.com/openziti/foundation/transport"
)

type Config struct {
	Datamesh *datamesh.Config `cf:"+required"`
	Peers    []transport.Address
}
