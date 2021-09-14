package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/openziti/foundation/transport"
)

type ProxyListener struct {
	bindAddress transport.Address
	txer EndpointTxer
	rxq  chan *dilithium.Buffer
}

func NewProxyListener(bindAddress transport.Address) *ProxyListener {
	return &ProxyListener{bindAddress: bindAddress}
}

func (pxl *ProxyListener) Connect(txer EndpointTxer, rxq chan *dilithium.Buffer) error {
	pxl.txer = txer
	pxl.rxq = rxq
}

type ProxyTerminator struct {
	dialAddress transport.Address
	txer EndpointTxer
	rxq  chan *dilithium.Buffer
}

func NewProxyTerminator(dialAddress transport.Address) *ProxyTerminator {
	return &ProxyTerminator{dialAddress: dialAddress}
}

func (pxt *ProxyTerminator) Connect(txer EndpointTxer, rxq chan *dilithium.Buffer) error {
	pxt.txer = txer
	pxt.rxq = rxq
}
