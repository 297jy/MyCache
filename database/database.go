package database

import "sync/atomic"

type CacheServer struct {
	dbSet []*atomic.Value
}

func NewStandaloneCacheServer() *CacheServer {
	return nil
}
