package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	set(key Key, value interface{}) bool
	get(key Key) (interface{}, bool)
	clear()
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
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) set(key Key, value interface{}) bool {
	_, ok := c.items[key]

	v := lruCacheItem{key: key, value: value}
	q := c.queue.PushFront(v)
	c.items[key] = q

	if c.queue.Len() > c.capacity {
		last := c.queue.Back()
		delete(c.items, last.Value.(lruCacheItem).key)
		c.queue.Remove(last)
	}
	return ok
}

func (c *lruCache) get(key Key) (interface{}, bool) {
	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.Remove(item)
	c.queue.PushFront(item.Value)
	return item.Value.(lruCacheItem).value, true
}

func (c *lruCache) clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	return c.set(key, value)
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	return c.get(key)
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.clear()
}
