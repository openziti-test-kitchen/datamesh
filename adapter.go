package datamesh

import (
	"github.com/openziti-incubator/datamesh/channel"
	"io"
)

type ChannelAdapter struct {
	ch channel.Channel
}

func NewChannelAdapter(ch channel.Channel) *ChannelAdapter {
	return &ChannelAdapter{ch}
}

func (ca *ChannelAdapter) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (ca *ChannelAdapter) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (ca *ChannelAdapter) Close() error {
	return nil
}