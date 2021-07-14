package datamesh

import (
	"github.com/openziti/foundation/channel2"
	"github.com/pkg/errors"
)

/*
 * Constants
 */

type ContentType uint32

const (
	ControlContentType ContentType = 10099
	PayloadContentType ContentType = 10100
	AckContentType     ContentType = 10101
)

type HeaderKey uint32

const (
	MinHeaderKey          = 2000
	ControlFlagsHeaderKey = 2256
)

type ControlFlag uint32

const (
	PingRequestControlFlag  ControlFlag = 1
	PingResponseControlFlag ControlFlag = 2
)

type PayloadFlag uint32

const (
	StartSessionPayloadFlag PayloadFlag = 1
	EndSessionPayloadFlag   PayloadFlag = 2
)

/*
 * Structures
 */

type Control struct {
	Flags uint32
	Data  []byte
}

type Routable struct {
	SessionId string
	Flags     uint32
}

type Payload struct {
	Routable
	Sequence int32
	Headers  map[uint8][]byte
	Data     []byte
}

type Acknowledgement struct {
	Routable
	Sequences []int32
}

/*
 * Implementation
 */

func (self *Control) Marshal() *channel2.Message {
	msg := channel2.NewMessage(int32(ControlContentType), self.Data)
	msg.PutUint32Header(ControlFlagsHeaderKey, self.Flags)
	return msg
}

func UnmarshallControl(msg *channel2.Message) (*Control, error) {
	control := &Control{Data: msg.Body}
	if flags, found := msg.GetUint32Header(ControlFlagsHeaderKey); found {
		control.Flags = flags
	} else {
		return nil, errors.Errorf("control message missing flags")
	}
	return control, nil
}
