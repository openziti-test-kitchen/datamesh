package datamesh

const (
	PayloadContentType = 10099
	AckContentType     = 10100
)

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