package hw04lrucache

import (
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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	item, ok := c.items[key]

	if !ok && c.queue.Len() == c.capacity {
		last := c.queue.Back()
		c.queue.Remove(last)
		for k, v := range c.items {
			if v == last {
				delete(c.items, k)
				break
			}
		}
	}

	switch {
	case !ok:
		item = c.queue.PushFront(value)
		c.items[key] = item
	default:
		item.Value = value
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

	return item.Value, true
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
