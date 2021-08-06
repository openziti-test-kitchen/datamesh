package datamesh

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequencesToAcks(t *testing.T) {
	seqs := []int32{0}
	ack := NewAcknowledgement("", seqs)
	arr := ack.sequencesToAcks()
	assert.Equal(t, 1, len(arr))
	assert.Equal(t, int32(0), arr[0].Start)
	assert.Equal(t, int32(0), arr[0].End)

	seqs = []int32{0, 1, 2, 3, 4}
	ack = NewAcknowledgement("", seqs)
	arr = ack.sequencesToAcks()
	assert.Equal(t, 1, len(arr))
	assert.Equal(t, int32(0), arr[0].Start)
	assert.Equal(t, int32(4), arr[0].End)

	seqs = []int32{0, 1, 2, 5, 6}
	ack = NewAcknowledgement("", seqs)
	arr = ack.sequencesToAcks()
	assert.Equal(t, 2, len(arr))
	assert.Equal(t, int32(0), arr[0].Start)
	assert.Equal(t, int32(2), arr[0].End)
	assert.Equal(t, int32(5), arr[1].Start)
	assert.Equal(t, int32(6), arr[1].End)

	seqs = []int32{}
	ack = NewAcknowledgement("", seqs)
	arr = ack.sequencesToAcks()
	assert.Nil(t, arr)  
}