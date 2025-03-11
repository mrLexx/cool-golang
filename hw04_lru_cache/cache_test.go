package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic, overflow and clean", func(t *testing.T) {
		c := NewCache(5)
		c.Set("a01", 100) // a01
		c.Set("a02", 200) // a02, a01
		c.Set("a03", 300) // a03, a02, a01
		c.Set("a04", 400) // a04, a03, a02, a01
		c.Set("a05", 500) // a05, a04, a03, a02, a01
		c.Set("a06", 500) // a06, a05, a04, a03, a02

		val, ok := c.Get("a01")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear()

		check := []Key{"a02", "a03", "a04", "a05", "a06"}
		for _, k := range check {
			val, ok = c.Get(k)
			require.False(t, ok)
			require.Nil(t, val)
		}
	})

	t.Run("purge logic, move to front by Set", func(t *testing.T) {
		c := NewCache(5)

		c.Set("a01", 100) // a01
		c.Set("a02", 200) // a02, a01
		c.Set("a03", 300) // a03, a02, a01
		c.Set("a04", 400) // a04, a03, a02, a01
		c.Set("a05", 500) // a05, a04, a03, a02, a01
		c.Set("a01", 100) // a01, a05, a04, a03, a02
		c.Set("a02", 200) // a02, a01, a05, a04, a03
		c.Set("a06", 600) // a06, a02, a01, a05, a04

		val, ok := c.Get("a03")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic, move to front by Get", func(t *testing.T) {
		c := NewCache(5)

		c.Clear()
		c.Set("a01", 100) // a01
		c.Set("a02", 200) // a02, a01
		c.Set("a03", 300) // a03, a02, a01
		c.Set("a04", 400) // a04, a03, a02, a01
		c.Set("a05", 500) // a05, a04, a03, a02, a01
		c.Get("a03")      // a03, a05, a04, a02, a01
		c.Get("a04")      // a04, a03, a05, a02, a01
		c.Get("a03")      // a03, a04, a05, a02, a01
		c.Set("a06", 600) // a06, a03, a04, a05, a02
		c.Set("a07", 700) // a07, a06, a03, a04, a05
		c.Set("a08", 800) // a08, a07, a06, a03, a04
		c.Set("a09", 900) // a09, a08, a07, a06, a03

		check := []Key{"a01", "a02", "a04", "a05"}
		for _, k := range check {
			val, ok := c.Get(k)
			require.False(t, ok)
			require.Nil(t, val)
		}

		check = []Key{"a03", "a06", "a07", "a08", "a09"}
		for _, k := range check {
			_, ok := c.Get(k)
			require.True(t, ok)
		}

		c.Set("a10", 1000)
		val, ok := c.Get("a01")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("capacity, less than zero", func(t *testing.T) {
		defer func() {
			r := recover()
			require.NotNil(t, r)
		}()

		c := NewCache(-1)
		c.Clear()
		_, ok := c.Get("a01")
		require.False(t, ok)

		c.Set("a01", 100) // panic
	})

	t.Run("capacity, equals zero", func(t *testing.T) {
		defer func() {
			r := recover()
			require.Nil(t, r)
		}()

		c := NewCache(0)
		c.Clear()
		c.Set("a01", 100) // []
		val, ok := c.Get("a01")
		require.Nil(t, val)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	_ = t

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
