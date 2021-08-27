package datamesh

import (
	"github.com/openziti/dilithium"
	"github.com/pkg/errors"
)

type Endpoint interface {
	Rx([]byte) error
}

type EndpointTxer interface {
	Tx(data []byte) error
}

type NIC interface {
	Destination
}

type nicImpl struct {
	address  Address
	endpoint Endpoint
	txp      *dilithium.TxPortal
	rxp      *dilithium.RxPortal
}

func newNIC(address Address, endpoint Endpoint) NIC {
	return &nicImpl{
		address:  address,
		endpoint: endpoint,
	}
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
