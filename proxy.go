package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type ProxyListener struct {
	bindAddress transport.Address
	in          chan transport.Connection
	conn        transport.Connection
	txq         EndpointTxer
	rxq         chan *dilithium.Buffer
	readBuf     []byte
}

func (pxl *ProxyListener) Accept(data []byte) error {
	_, err := pxl.conn.Writer().Write(data)
	return err
}

func (pxl *ProxyListener) Close() {
	if err := pxl.conn.Close(); err != nil {
		logrus.WithError(err).Errorf("failure closing proxy listener conn")
	}
}

func NewProxyListener(bindAddress transport.Address) *ProxyListener {
	return &ProxyListener{bindAddress: bindAddress, in: make(chan transport.Connection, 1), readBuf: make([]byte, 128*1024)}
}

func (pxl *ProxyListener) Connect(txq EndpointTxer) error {
	_, err := pxl.bindAddress.Listen("ProxyListener", nil, pxl.in, nil)
	if err != nil {
		return errors.Wrap(err, "error listening")
	}
	go pxl.accept()

	pxl.txq = txq

	return nil
}

func (pxl *ProxyListener) accept() {
	logrus.Infof("listening [%v]", pxl.bindAddress)
	select {
	case conn := <-pxl.in:
		pxl.conn = conn
	}
	logrus.Infof("accepted connection [%v]", pxl.conn.Detail())
	go pxl.txer()
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

type ProxyListenerFactory struct {
	BindAddress transport.Address
	CircuitId   Circuit
}

func (pxlf *ProxyListenerFactory) Create() (Endpoint, error) {
	return NewProxyListener(pxlf.BindAddress), nil
}

func (pxlf *ProxyListenerFactory) Circuit() Circuit {
	return pxlf.CircuitId
}

type ProxyTerminator struct {
	dialAddress transport.Address
	conn        transport.Connection
	txq         EndpointTxer
	rxq         chan *dilithium.Buffer
	readBuf     []byte
}

func NewProxyTerminator(dialAddress transport.Address) *ProxyTerminator {
	return &ProxyTerminator{dialAddress: dialAddress, readBuf: make([]byte, 128*1024)}
}

func (pxt *ProxyTerminator) Accept(data []byte) error {
	_, err := pxt.conn.Writer().Write(data)
	return err
}

func (pxt *ProxyTerminator) Close() {
	if err := pxt.conn.Close(); err != nil {
		logrus.WithError(err).Errorf("failure closing proxy terminator conn")
	}
}
func (pxt *ProxyTerminator) Connect(txq EndpointTxer) error {
	conn, err := pxt.dialAddress.Dial("ProxyTerminator", nil, 5*time.Second, nil)
	if err != nil {
		return errors.Wrap(err, "error dialing")
	}
	logrus.Infof("connection dialed [%v]", pxt.dialAddress)
	pxt.conn = conn

	pxt.txq = txq
	go pxt.txer()

	return nil
}

func (pxt *ProxyTerminator) txer() {
	logrus.Info("started")
	defer logrus.Info("exited")

	for {
		if n, err := pxt.conn.Reader().Read(pxt.readBuf); err == nil {
			if err := pxt.txq.Tx(pxt.readBuf[:n]); err != nil {
				logrus.Errorf("forward error (%v)", err)
				return
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

type ProxyTerminatorFactory struct {
	DialAddress transport.Address
	CircuitId   Circuit
}

func (pxtf *ProxyTerminatorFactory) Create() (Endpoint, error) {
	return NewProxyTerminator(pxtf.DialAddress), nil
}

func (pxtf *ProxyTerminatorFactory) Circuit() Circuit {
	return pxtf.CircuitId
}

type ProxyFactory interface {
	Create() (Endpoint, error)
	Circuit() Circuit
}
