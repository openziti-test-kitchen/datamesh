package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/openziti/dilithium/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Endpoint defines the primary "extensible" component in datamesh. An Endpoint sits inside of a NIC, which allows it to
// communicate with another NIC, and its contained Endpoint elsewhere on the network.
//
type Endpoint interface {
	Connect(txer EndpointTxer, rxer chan *dilithium.Buffer) error
}

// EndpointTxer defines the transmitter interface exposed to an Endpoint.
//
type EndpointTxer interface {
	ToNetwork(data []byte) error
}

type NIC interface {
	Destination
}

type nicImpl struct {
	circuit  Circuit
	address  Address
	endpoint Endpoint
	dm       *Datamesh
	da       *NICAdapter
	seq      *util.Sequence
	txa      dilithium.TxAlgorithm
	txp      *dilithium.TxPortal
	rxp      *dilithium.RxPortal
	rxq      chan *dilithium.Buffer
	closer   *dilithium.Closer
	pool     *dilithium.Pool
}

func newNIC(dm *Datamesh, circuit Circuit, address Address, endpoint Endpoint) NIC {
	nic := &nicImpl{
		circuit:  circuit,
		address:  address,
		endpoint: endpoint,
		dm:       dm,
		seq:      util.NewSequence(0),
		pool:     dilithium.NewPool("nic", 128*1024),
	}
	nic.da = NewNICAdapter(nic)
	return nic
}

func (nic *nicImpl) SetTxAlgorithm(txa dilithium.TxAlgorithm) error {
	if nic.txa != nil {
		return errors.New("algorithm already present")
	}
	nic.txa = txa
	return nil
}

func (nic *nicImpl) Start() error {
	if nic.closer == nil && nic.txp == nil && nic.rxp == nil {
		nic.closer = dilithium.NewCloser(nic.seq, nil)
		nic.txp = dilithium.NewTxPortal(nic.da, nic.txa, nic.closer)
		nic.rxp = dilithium.NewRxPortal(nic.da, nic.txp, nic.seq, nic.closer)
		nic.txp.Start()
		go nic.rxer()
		if err := nic.endpoint.Connect(nic, nic.rxq); err != nil {
			return errors.Wrap(err, "unable to start nic")
		}
		logrus.Info("started")
		return nil

	} else {
		return errors.New("already started")
	}
}

func (nic *nicImpl) Address() Address {
	return nic.address
}

func (nic *nicImpl) FromNetwork(data []byte) error {
	return errors.Errorf("not implemented")
}

func (nic *nicImpl) Close() error {
	return errors.Errorf("not implemented")
}

func (nic *nicImpl) ToNetwork(data []byte) error {
	return errors.Errorf("not implemented")
}

func (nic *nicImpl) rxer() {
	logrus.Info("started")
	defer logrus.Info("exited")

	for {
		buf := nic.pool.Get()
		n, err := nic.rxp.Read(buf.Data)
		if err != nil {
			logrus.Errorf("read error (%v)", err)
		}
		buf.Used = uint32(n)
		select {
		case nic.rxq <- buf:
		default:
			logrus.Info("dropped")
		}
	}
}
