package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/openziti/foundation/transport"
	"github.com/sirupsen/logrus"
	"io"
)

type ProxyListener struct {
	bindAddress transport.Address
	conn        transport.Connection
	txq         EndpointTxer
	rxq         chan *dilithium.Buffer
	readBuf     []byte
}

func NewProxyListener(bindAddress transport.Address) *ProxyListener {
	return &ProxyListener{bindAddress: bindAddress}
}

func (pxl *ProxyListener) Connect(txq EndpointTxer, rxq chan *dilithium.Buffer) error {
	// listen for 1 connection

	pxl.txq = txq
	pxl.rxq = rxq
	go pxl.rxer()
	go pxl.txer()

	return nil
}

func (pxl *ProxyListener) rxer() {
	logrus.Info("started")
	defer logrus.Info("exited")

	for {
		select {
		case buf := <-pxl.rxq:
			if n, err := pxl.conn.Writer().Write(buf.Data[:buf.Used]); err == nil {
				if uint32(n) != buf.Used {
					logrus.Warn("short write")
				}
			} else {
				logrus.Errorf("write error (%v)", err)
			}
			buf.Unref()
		}
	}
}

func (pxl *ProxyListener) txer() {
	logrus.Info("started")
	defer logrus.Info("exited")

	for {
		if n, err := pxl.conn.Reader().Read(pxl.readBuf); err == nil {
			if err := pxl.txq.Tx(pxl.readBuf[:n]); err != nil {
				logrus.Errorf("forward error (%v)", err)
			}
		} else if err == io.EOF {
			// close handling
			logrus.Warn("EOF")
			return
		} else {
			logrus.Errorf("read error (%v)", err)
		}
	}
}

type ProxyTerminator struct {
	dialAddress transport.Address
	conn        transport.Connection
	txq         EndpointTxer
	rxq         chan *dilithium.Buffer
	readBuf     []byte
}

func NewProxyTerminator(dialAddress transport.Address) *ProxyTerminator {
	return &ProxyTerminator{dialAddress: dialAddress}
}

func (pxt *ProxyTerminator) Connect(txq EndpointTxer, rxq chan *dilithium.Buffer) error {
	// dial service

	pxt.txq = txq
	pxt.rxq = rxq
	go pxt.rxer()
	go pxt.txer()

	return nil
}

func (pxt *ProxyTerminator) rxer() {
	logrus.Info("started")
	defer logrus.Info("exited")

	for {
		select {
		case buf := <-pxt.rxq:
			if n, err := pxt.conn.Writer().Write(buf.Data[:buf.Used]); err == nil {
				if uint32(n) != buf.Used {
					logrus.Warn("short write")
				}
			} else {
				logrus.Errorf("write error (%v)", err)
			}
			buf.Unref()
		}
	}
}

func (pxt *ProxyTerminator) txer() {
	logrus.Info("started")
	defer logrus.Info("exited")

	for {
		if n, err := pxt.conn.Reader().Read(pxt.readBuf); err == nil {
			if err := pxt.txq.Tx(pxt.readBuf[:n]); err != nil {
				logrus.Errorf("forward error (%v)", err)
			}
		} else if err == io.EOF {
			// close handling
			logrus.Warn("EOF")
			return
		} else {
			logrus.Errorf("read error (%v)", err)
		}
	}
}
