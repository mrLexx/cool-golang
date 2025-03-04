package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type lruCacheItem struct {
	k Key
	v interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	_, ok := c.items[key]

	if !ok && c.queue.Len() >= c.capacity {
		last := c.queue.Back()

		delete(c.items, last.Value.(lruCacheItem).k)
		c.queue.Remove(last)
	}

	v := lruCacheItem{k: key, v: value}
	elem := c.queue.PushFront(v)
	c.items[key] = elem

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	elem, ok := c.items[key]

	if !ok {
		return nil, false
	}
	c.queue.PushFront(elem)

	return elem.Value.(lruCacheItem).v, true
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue.Clear()
}
