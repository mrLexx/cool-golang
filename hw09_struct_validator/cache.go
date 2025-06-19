package hw09structvalidator

import (
	"sync"
)

type cacheValue struct {
	Val  string
	Rule string
}
type cache[T any] interface {
	Set(val cacheValue, result T) bool
	Get(val cacheValue) (T, bool)
	Reset()
}

type cacheItem[T any] struct {
	lock  sync.RWMutex
	Items map[cacheValue]T
}

func NewCache[T any]() cache[T] {
	return &cacheItem[T]{
		Items: make(map[cacheValue]T),
	}
}

func (l *cacheItem[T]) Set(val cacheValue, value T) bool {
	l.lock.Lock()
	l.Items[val] = value
	l.lock.Unlock()
	_, ok := l.Get(val)
	return ok
}

func (l *cacheItem[T]) Get(val cacheValue) (T, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	v, ok := l.Items[val]
	return v, ok
}

func (l *cacheItem[T]) Reset() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Items = make(map[cacheValue]T)
}
