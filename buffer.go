package datamesh

import (
	"sync"
	"sync/atomic"
)

type Buffer struct {
	data []byte
	refs int32
	uz   uint32
	sz   uint32
	pool *Pool
}

func NewBuffer(pool *Pool) *Buffer {
	return &Buffer{
		data: make([]byte, pool.bufferSz),
		pool: pool,
		uz:   0,
		sz:   pool.bufferSz,
		refs: 0,
	}
}

func (self *Buffer) Ref() {
	atomic.AddInt32(&self.refs, 1)
}

func (self *Buffer) Unref() {
	if atomic.AddInt32(&self.refs, -1) < 1 {
		self.uz = 0
		self.pool.put(self)
	}
}

type Pool struct {
	id       string
	bufferSz uint32
	store    *sync.Pool
}

func NewPool(id string, bufferSz uint32) *Pool {
	pool := &Pool{
		id:       id,
		bufferSz: bufferSz,
		store:    new(sync.Pool),
	}
	pool.store.New = pool.allocate
	return pool
}

func (self *Pool) Get() *Buffer {
	buf := self.store.Get().(*Buffer)
	buf.Ref()
	return buf
}

func (self *Pool) put(buffer *Buffer) {
	self.store.Put(buffer)
}

func (self *Pool) allocate() interface{} {
	return NewBuffer(self)
}
