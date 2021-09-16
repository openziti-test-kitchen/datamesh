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
	conn        transport.Connection
	txq         EndpointTxer
	rxq         chan *dilithium.Buffer
	readBuf     []byte
}

func NewProxyListener(bindAddress transport.Address) *ProxyListener {
	return &ProxyListener{bindAddress: bindAddress, readBuf: make([]byte, 128 * 1024)}
}

func (pxl *ProxyListener) Connect(txq EndpointTxer, rxq chan *dilithium.Buffer) error {
	in := make(chan transport.Connection)
	_, err := pxl.bindAddress.Listen("ProxyListener", nil, in, nil)
	if err != nil {
		return errors.Wrap(err, "error listening")
	}
	logrus.Infof("listening [%v]", pxl.bindAddress)
	select {
	case conn := <-in:
		pxl.conn = conn
	}
	logrus.Infof("accepted connection [%v]", pxl.conn.Detail())

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
	return &ProxyTerminator{dialAddress: dialAddress, readBuf: make([]byte, 128 * 1024)}
}

func (pxt *ProxyTerminator) Connect(txq EndpointTxer, rxq chan *dilithium.Buffer) error {
	conn, err := pxt.dialAddress.Dial("ProxyTerminator", nil, 5*time.Second, nil)
	if err != nil {
		return errors.Wrap(err, "error dialing")
	}
	logrus.Infof("connection dialed [%v]", pxt.dialAddress)
	pxt.conn = conn

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
				return
			}
			logrus.Infof("tx (%v bytes)", n)
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
