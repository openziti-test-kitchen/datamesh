package datamesh

import "github.com/pkg/errors"

type NICTxer interface {
	Tx(data []byte)
}

type NIC interface {
	Destination
}

type nicImpl struct {
	address  Address
	endpoint Endpoint
}

func newNIC(address Address, endpoint Endpoint) NIC {
	return &nicImpl{address, endpoint}
}

func (self *nicImpl) Address() Address {
	return self.address
}

func (self *nicImpl) SendData(data *Data) error {
	return errors.Errorf("not implemented")
}

func (self *nicImpl) Close() error {
	return errors.Errorf("not implemented")
}