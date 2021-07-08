package datamesh

import (
	"github.com/openziti/foundation/channel2"
	"github.com/sirupsen/logrus"
)

type bindHandler struct{}

func (_ *bindHandler) BindChannel(ch channel2.Channel) error {
	ch.AddReceiveHandler(&payloadReceiveHandler{})
	ch.AddReceiveHandler(&ackReceiveHandler{})
	return nil
}

type payloadReceiveHandler struct{}

func (_ *payloadReceiveHandler) ContentType() int32 {
	return int32(PayloadContentType)
}

func (_ *payloadReceiveHandler) HandleReceive(m *channel2.Message, _ channel2.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}

type ackReceiveHandler struct{}

func (_ *ackReceiveHandler) ContentType() int32 {
	return int32(AckContentType)
}

func (_ *ackReceiveHandler) HandleReceive(m *channel2.Message, _ channel2.Channel) {
	logrus.Infof("received [%d] bytes", len(m.Body))
}
