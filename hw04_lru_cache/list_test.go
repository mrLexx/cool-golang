package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("front and back", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // 10
		require.Equal(t, true, l.Front() == l.Back())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})

	t.Run("move to front", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)

		front := l.Front()
		l.MoveToFront(front)
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)

		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)
		l.PushBack(50) // 10, 20, 30, 40, 50

		front = l.Front()
		l.MoveToFront(front) // [10, 20, 30, 40, 50]
		require.Equal(t, 10, l.Front().Value)

		back := l.Back()
		l.MoveToFront(back) // [50, 10, 20, 30, 40]
		require.Equal(t, 50, l.Front().Value)
	})

	t.Run("remove", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)
		l.PushBack(50) // 10, 20, 30, 40, 50

		back := l.Back()
		l.Remove(back) // [10, 20, 30, 40]
		require.Equal(t, 40, l.Back().Value)
		require.Equal(t, 4, l.Len())

		front := l.Front()
		l.Remove(front) // [20, 30, 40]
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 3, l.Len())

		l.Remove(l.Front()) // [30, 40]
		l.Remove(l.Front()) // [40]
		l.Remove(l.Front()) // []
		l.Remove(l.Front()) // []
		require.Equal(t, 0, l.Len())

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{}, elems)
	})
}
