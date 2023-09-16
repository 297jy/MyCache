package cache

import (
	"container/list"
	"fmt"
	"strconv"
)

type lruCache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
}

func (c *lruCache) Put(key string, value ByteView) (ok bool) {
	valueBytes := int64(len(key)) + int64(value.Len())
	// 极端情况，最大的容量都装不下新的值直接报错
	if c.maxBytes < valueBytes {
		return false
	}

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += valueBytes
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.Eliminate()
	}
	return true
}

func (c *lruCache) Eliminate() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		freeBytes := int64(len(kv.key)) + int64(kv.value.Len())
		c.nbytes -= freeBytes

		fmt.Println("淘汰的key为:" + kv.key + "," + strconv.FormatInt(freeBytes, 10))
	}
}

func (c *lruCache) Get(key string) (value ByteView, ok bool) {

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *lruCache) Remove(key string) (value ByteView, ok bool) {

	if ele, ok := c.cache[key]; ok {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		freeBytes := int64(len(kv.key)) + int64(kv.value.Len())
		c.nbytes -= freeBytes
		return kv.value, true
	}
	return value, false
}

func NewLruCache(maxBytes int64) Cache {
	return &lruCache{
		maxBytes: maxBytes,
		nbytes:   0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}
