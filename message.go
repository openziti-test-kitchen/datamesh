package datamesh

const (
	ControlContentType = 10099
	PayloadContentType = 10100
	AckContentType     = 10101
)

type Control struct {
	Headers map[uint8][]byte
}

type Header struct {
	SessionId string
	Flags     uint32
}

type Payload struct {
	Header
	Sequence int32
	Headers  map[uint8][]byte
	Data     []byte
}

type Acknowledgement struct {
	Header
	Sequences []int32
}
