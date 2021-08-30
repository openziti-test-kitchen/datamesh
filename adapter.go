package datamesh

import (
	"io"
)

type DestinationAdapter struct {
	dm      *Datamesh
	circuit CircuitId
	srcAddr Address
	dstAddr Address
}

func NewDestinationAdapter(dm *Datamesh, circuit CircuitId, srcAddr, dstAddr Address) *DestinationAdapter {
	return &DestinationAdapter{dm, circuit, srcAddr, dstAddr}
}

func (ca *DestinationAdapter) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (ca *DestinationAdapter) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (ca *DestinationAdapter) Close() error {
	return nil
}
