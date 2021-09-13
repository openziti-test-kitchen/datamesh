package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/openziti/foundation/transport"
)

type ProxyListener struct {
	bindAddress transport.Address
	txer        EndpointTxer
	rxer        chan *dilithium.Buffer
}

func (pxl *ProxyListener) Connect(txer EndpointTxer, rxer chan *dilithium.Buffer) {
	pxl.txer = txer
	pxl.rxer = rxer
}

type ProxyTerminator struct {
	dialAddress transport.Address
	txer        EndpointTxer
	rxer        chan *dilithium.Buffer
}

func (pxt *ProxyTerminator) Connect(txer EndpointTxer, rxer chan *dilithium.Buffer) {
	pxt.txer = txer
	pxt.rxer = rxer
}
