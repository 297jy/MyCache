package pool

import (
	"errors"
	"sync"
)

var (
	ErrClosed = errors.New("pool closed")
	ErrMax    = errors.New("reach max connection limit")
)

type request chan interface{}

type Config struct {
	MaxIdle   uint
	MaxActive uint
}

// Pool 请求连接池
type Pool struct {
	Config
	factory     func() (interface{}, error)
	finalizer   func(x interface{})
	idles       chan interface{}
	waitingReqs []request
	activeCount uint // increases during creating connection, decrease during destroying connection
	mu          sync.Mutex
	closed      bool
}

func New(factory func() (interface{}, error), finalizer func(x interface{}), cfg Config) *Pool {
	return &Pool{
		factory:     factory,
		finalizer:   finalizer,
		idles:       make(chan interface{}, cfg.MaxIdle),
		waitingReqs: make([]request, 0),
		Config:      cfg,
	}
}

func (pool *Pool) Get() (interface{}, error) {
	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()
		return nil, ErrClosed
	}

	select {
	case item := <-pool.idles:
		pool.mu.Unlock()
		return item, nil
	default:
		// no pooled item, create one
		return pool.getOnNoIdle()
	}
}

func (pool *Pool) getOnNoIdle() (interface{}, error) {
	if pool.activeCount >= pool.MaxActive {
		// waiting for connection being returned
		req := make(chan interface{}, 1)
		pool.waitingReqs = append(pool.waitingReqs, req)
		pool.mu.Unlock()
		x, ok := <-req
		if !ok {
			return nil, ErrMax
		}
		return x, nil
	}

	// create a new connection
	pool.activeCount++ // hold a place for new connection
	pool.mu.Unlock()
	x, err := pool.factory()
	if err != nil {
		// create failed return token
		pool.mu.Lock()
		pool.activeCount-- // release the holding place
		pool.mu.Unlock()
		return nil, err
	}
	return x, nil
}
