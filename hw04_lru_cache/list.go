package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len int

	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (list list) Len() int {
	return list.len
}

func (list list) Front() *ListItem {
	return list.front
}

func (list list) Back() *ListItem {
	return list.back
}

func (list *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: nil, Prev: nil}

	switch list.len {
	case 0:
		list.front, list.back = item, item
	default:
		list.front.Prev = item
		item.Next = list.front
		list.front = item
	}

	list.len++
	return list.front
}

func (list *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: nil, Prev: nil}

	switch list.len {
	case 0:
		list.front, list.back = item, item

	default:
		list.back.Next = item
		item.Prev = list.back
		list.back = item
	}

	list.len++
	return list.front
}

func (list *list) Remove(i *ListItem) {
	if list.len == 0 {
		return
	}

	switch {
	case i.Prev == nil && i.Next == nil:
		list.back, list.front = nil, nil

	case i.Prev != nil && i.Next != nil:
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev

	case i.Prev != nil:
		i.Prev.Next = nil
		list.back = i.Prev

	case i.Next != nil:
		i.Next.Prev = nil
		list.front = i.Next
	}

	list.len--
}

func (list *list) MoveToFront(i *ListItem) {
	if list.front == i {
		return
	}

	switch {
	case i.Prev != nil && i.Next != nil:
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev

	case i.Prev != nil:
		i.Prev.Next = nil
		list.back = i.Prev
	}
	i.Prev, i.Next = nil, list.front
	list.front.Prev, list.front = i, i
}
