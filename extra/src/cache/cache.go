package cache

import (
	"sync"
	"unsafe"
)

type Cache interface {
	Put(key string, value ByteView) (ok bool)
	Get(key string) (value ByteView, ok bool)
	Remove(key string) (value ByteView, ok bool)
}

type Value interface {
	Len() int
}

type ByteView struct {
	b       []byte
	version int
}

func (v ByteView) Len() int {
	return len(v.b) + int(unsafe.Sizeof(v.version))
}

func (v ByteView) Version() int {
	return v.version
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

type entry struct {
	key   string
	value ByteView
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache Cache
}

type NewCacheFunc func(maxBytes int64) Cache

type NewCacheBuilder struct {
	maxBytes    int64
	builderFunc NewCacheFunc
}

type GroupManage interface {
	Put(groupName string, key string, view ByteView) (ok bool)
	Get(groupName string, key string) (byteView ByteView, ok bool)
	Remove(groupName string, key string) (value ByteView, ok bool)
}

type groupContainer struct {
	mu           sync.RWMutex
	cacheBuilder NewCacheBuilder
	groups       map[string]*Group
}

func (g *groupContainer) Put(groupName string, key string, byteView ByteView) (ok bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	group, ok := g.groups[groupName]
	if !ok {
		group = &Group{
			name:      groupName,
			mainCache: g.cacheBuilder.builderFunc(g.cacheBuilder.maxBytes),
		}
		g.groups[groupName] = group
	}
	return group.mainCache.Put(key, byteView)
}

func (g *groupContainer) Get(groupName string, key string) (byteView ByteView, ok bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	group, ok := g.groups[groupName]
	if !ok {
		return
	}
	return group.mainCache.Get(key)
}

func (g *groupContainer) Remove(groupName string, key string) (value ByteView, ok bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	group, ok := g.groups[groupName]
	if !ok {
		return value, false
	}
	return group.mainCache.Remove(key)
}

func NewGroupManage(newCacheBuilder NewCacheBuilder) GroupManage {
	return &groupContainer{
		cacheBuilder: newCacheBuilder,
		groups:       make(map[string]*Group),
	}
}
