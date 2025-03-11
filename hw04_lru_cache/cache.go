package hw04lrucache

import (
	"errors"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	queue    List
	items    map[Key]*ListItem
	capacity int
}

type lruCacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	if capacity < 0 {
		err := errors.New("capacity must >= 0")
		panic(err)
	}

	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	if c.capacity == 0 {
		return false
	}

	item, ok := c.items[key]

	if !ok && c.queue.Len() == c.capacity {
		last := c.queue.Back()
		c.queue.Remove(last)
		delete(c.items, last.Value.(lruCacheItem).key)
	}

	v := lruCacheItem{key: key, value: value}
	switch {
	case !ok:
		item = c.queue.PushFront(v)
		c.items[key] = item
	default:
		item.Value = v
		c.queue.MoveToFront(item)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return item.Value.(lruCacheItem).value, true
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
