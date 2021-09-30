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
	Connect(txer EndpointTxer) error
	dilithium.Sink
}

// EndpointTxer defines the transmitter interface exposed to an Endpoint.
//
type EndpointTxer interface {
	Tx(data []byte) error
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
	netq     chan *dilithium.Buffer
	closer   *dilithium.Closer
	pool     *dilithium.Pool
	ii       dilithium.InstrumentInstance
}

func newNIC(dm *Datamesh, circuit Circuit, address Address, endpoint Endpoint, ii dilithium.InstrumentInstance) NIC {
	logrus.Info("started")
	nic := &nicImpl{
		circuit:  circuit,
		address:  address,
		endpoint: endpoint,
		dm:       dm,
		seq:      util.NewSequence(0),
		netq:     make(chan *dilithium.Buffer, 16),
		pool:     dilithium.NewPool("nic", 128*1024, ii),
		ii:       ii,
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
		nic.txp = dilithium.NewTxPortal(nic.da, nic.txa, nic.closer, nic.ii)
		nic.rxp = dilithium.NewRxPortal(nic.da, nic.endpoint, nic.txp, nic.seq, nic.closer, nic.ii)
		nic.txp.Start()
		if err := nic.endpoint.Connect(nic); err != nil {
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

func (nic *nicImpl) FromNetwork(payload *Payload) error {
	buf := nic.pool.Get()
	n := copy(buf.Data, payload.Buf.Data[:payload.Buf.Used])
	buf.Used = uint32(n)
	nic.netq <- buf
	return nil
}

func (nic *nicImpl) Close() error {
	return errors.Errorf("not implemented")
}

func (nic *nicImpl) Tx(data []byte) error {
	n, err := nic.txp.Tx(data, nic.seq)
	if err != nil {
		return errors.Wrap(err, "to network")
	}
	if n != len(data) {
		return errors.New("short to network")
	}
	return nil
}
