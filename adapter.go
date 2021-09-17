package datamesh

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
)

type NICAdapter struct {
	nic NIC
}

func NewNICAdapter(nic NIC) *NICAdapter {
	return &NICAdapter{nic}
}

func (na *NICAdapter) Read(p []byte) (n int, err error) {
	select {
	case buf, ok := <-na.nic.(*nicImpl).netq:
		if ok {
			n := copy(p, buf.Data[:buf.Used])
			logrus.Infof("read (%v)", n)
			return n, nil
		} else {
			logrus.Warn("no read, not ok")
		}
	}
	return 0, io.EOF
}

func (na *NICAdapter) Write(p []byte) (n int, err error) {
	nic := na.nic.(*nicImpl)

	logrus.Infof("tx (%v, %v)", len(p), nic.address)

	payload := NewPayload(nic.circuit)
	payload.Buf = nic.pool.Get()
	n = copy(payload.Buf.Data, p)
	if n != len(p) {
		return 0, errors.New("short copy")
	}
	payload.Buf.Used = uint32(n)

	if err := nic.dm.Fwd.Forward(nic.address, payload); err != nil {
		return 0, errors.Wrap(err, "error forwarding")
	}

	return n, nil
}

func (na *NICAdapter) Close() error {
	return nil
}
