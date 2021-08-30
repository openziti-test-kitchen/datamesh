package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/openziti/dilithium/util"
	"github.com/pkg/errors"
)

// Endpoint defines the primary "extensible" component in datamesh. An Endpoint sits inside of a NIC, which allows it to
// communicate with another NIC, and its contained Endpoint elsewhere on the network.
//
type Endpoint interface {
	Connect(txer EndpointTxer, rxer chan []byte)
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
	closer   *dilithium.Closer
}

func newNIC(dm *Datamesh, circuit Circuit, address Address, endpoint Endpoint) NIC {
	nic := &nicImpl{
		circuit:  circuit,
		address:  address,
		endpoint: endpoint,
		dm:       dm,
		seq:      util.NewSequence(0),
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

func (nic *nicImpl) Start() {
}

func (nic *nicImpl) Address() Address {
	return nic.address
}

func (nic *nicImpl) SendData(data *Data) error {
	return errors.Errorf("not implemented")
}

func (nic *nicImpl) Close() error {
	return errors.Errorf("not implemented")
}

func (nic *nicImpl) Tx(data []byte) error {
	return errors.Errorf("not implemented")
}
