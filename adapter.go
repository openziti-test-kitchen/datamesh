package datamesh

import (
	"io"
)

type NICAdapter struct {
	nic NIC
}

func NewNICAdapter(nic NIC) *NICAdapter {
	return &NICAdapter{nic}
}

func (na *NICAdapter) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (na *NICAdapter) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (na *NICAdapter) Close() error {
	return nil
}
