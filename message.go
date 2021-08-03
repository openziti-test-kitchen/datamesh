package datamesh

import (
	"encoding/binary"
	"github.com/openziti-incubator/datamesh/channel"
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
	MinHeaderKey              = 2000
	ControlFlagsHeaderKey     = 2256
	PingIdHeaderKey           = 2257
	PingTimestampHeaderKey    = 2258
	PayloadSequenceHeaderKey  = 2300
	PayloadSessionIdHeaderKey = 2301
	PayloadFlagsHeaderKey     = 2302
	MaxHeaderKey              = 2999
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
	Flags   uint32
	Headers map[int32][]byte
}

type Payload struct {
	Sequence  int32
	SessionId string
	Flags     uint32
	Headers   map[int32][]byte
	Data      []byte
}

type Acknowledgement struct {
	SessionId string
	Sequences []int32
}

/*
 * Implementation
 */

func NewControl(flags uint32, headers map[int32][]byte) *Control {
	return &Control{flags, headers}
}

func (self *Control) Marshal() *channel.Message {
	msg := channel.NewMessage(int32(ControlContentType), nil)
	msg.PutUint32Header(ControlFlagsHeaderKey, self.Flags)
	for k, v := range self.Headers {
		msg.Headers[k] = v
	}
	return msg
}

func UnmarshalControl(msg *channel.Message) (*Control, error) {
	control := &Control{}
	if flags, found := msg.GetUint32Header(ControlFlagsHeaderKey); found {
		control.Flags = flags
	} else {
		return nil, errors.Errorf("control message missing flags")
	}
	for k, v := range msg.Headers {
		if k >= MinHeaderKey && k <= MaxHeaderKey {
			if control.Headers == nil {
				control.Headers = make(map[int32][]byte)
			}
			control.Headers[k] = v
		}
	}
	return control, nil
}

func NewPayload(sequence int32, sessionId string) *Payload {
	return &Payload{Sequence: sequence, SessionId: sessionId}
}

func (self *Payload) Marshal() *channel.Message {
	msg := channel.NewMessage(int32(PayloadContentType), self.Data)
	for k, v := range self.Headers {
		msg.Headers[k] = v
	}
	msg.PutUint32Header(PayloadSequenceHeaderKey, uint32(self.Sequence))
	msg.Headers[PayloadSessionIdHeaderKey] = []byte(self.SessionId)
	if self.Flags > 0 {
		msg.PutUint32Header(PayloadFlagsHeaderKey, self.Flags)
	}
	return msg
}

func UnmarshalPayload(msg *channel.Message) (*Payload, error) {
	p := &Payload{Data: make([]byte, len(msg.Body))}
	copy(p.Data, msg.Body)
	if sequence, found := msg.GetUint32Header(PayloadSequenceHeaderKey); found {
		p.Sequence = int32(sequence)
	} else {
		return nil, errors.New("missing sequence from payload")
	}
	if sessionId, found := msg.Headers[PayloadSessionIdHeaderKey]; found {
		p.SessionId = string(sessionId)
	} else {
		return nil, errors.New("missing sessionId from payload")
	}
	if flags, found := msg.GetUint32Header(PayloadFlagsHeaderKey); found {
		p.Flags = flags
	}
	for k, v := range msg.Headers {
		if k < MinHeaderKey && k > MaxHeaderKey {
			p.Headers[k] = v
		}
	}
	return p, nil
}

/*
 * Utils
 */

type headers map[int32][]byte

func newHeaders() headers {
	return make(map[int32][]byte)
}

func (self *headers) PutBytes(key int32, value []byte) {
	encoded := make([]byte, len(value))
	copy(encoded, value)
	map[int32][]byte(*self)[key] = encoded
}

func (self *headers) GetBytes(key int32) ([]byte, bool) {
	v, found := map[int32][]byte(*self)[key]
	return v, found
}

func (self *headers) PutInt64(key int32, value int64) {
	encoded := make([]byte, 8)
	binary.LittleEndian.PutUint64(encoded, uint64(value))
	map[int32][]byte(*self)[key] = encoded
}

func (self *headers) GetInt64(key int32) (int64, bool) {
	encoded, ok := map[int32][]byte(*self)[key]
	if !ok || len(encoded) != 8 {
		return 0, false
	}
	result := int64(binary.LittleEndian.Uint64(encoded))
	return result, true
}
