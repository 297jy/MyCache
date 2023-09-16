package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	cacheBuilder := NewCacheBuilder{
		builderFunc: NewLruCache,
		maxBytes:    40,
	}
	g := NewGroupManage(cacheBuilder)
	g.Put("test", "test", ByteView{
		b:       []byte("test"),
		version: 4,
	})
	g.Put("test", "test1", ByteView{
		b:       []byte("test"),
		version: 4,
	})
	g.Put("test", "test2", ByteView{
		b:       []byte("test"),
		version: 4,
	})

}
