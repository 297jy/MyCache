package connection

import "sync"

// FakeConn implements redis.Connection for test
type FakeConn struct {
	Connection
	buf    []byte
	offset int
	waitOn chan struct{}
	closed bool
	mu     sync.Mutex
}

func NewFakeConn() *FakeConn {
	c := &FakeConn{}
	return c
}
