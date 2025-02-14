package linkedlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/johannessarpola/gollections/internal/node"
)

type LinkedList[T comparable] struct {
	head *node.Node[T]
	mu   sync.Mutex
}

func NewLinkedList[T comparable]() LinkedList[T] {
	return LinkedList[T]{
		head: nil,
		mu:   sync.Mutex{},
	}
}

func (l *LinkedList[T]) withLock(f func()) {
	defer l.mu.Unlock()
	l.mu.Lock()
	f()
}

func (l *LinkedList[T]) String() string {
	sb := strings.Builder{}
	var s []string
	sb.WriteString("[")
	for _, v := range l.All {
		s = append(s, fmt.Sprintf("%v", v))
	}
	sb.WriteString(strings.Join(s, ","))
	sb.WriteString("]")
	return sb.String()
}

func (l *LinkedList[T]) Contains(value T) bool {
	b := false
	l.withLock(func() {
		for current := l.head; current != nil; current = current.Next {
			if current.Inner == value {
				b = true
				return
			}
		}
	})
	return b
}

func (l *LinkedList[T]) Append(value T) {
	n := node.NewNode(value)
	l.withLock(func() {
		if l.head == nil {
			l.head = &n
			return
		}

		current := l.head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = &n
	})
}

func join[T comparable](root *node.Node[T], another *node.Node[T]) *node.Node[T] {
	if root == nil {
		return another
	}

	last := root
	for current := root; current != nil; current = current.Next {
		last = current
	}

	last.Next = another
	return root
}

func (l *LinkedList[T]) Join(another *LinkedList[T]) *LinkedList[T] {
	l.withLock(func() {
		l.head = join(l.head, another.head)
	})

	return l
}

func (l *LinkedList[T]) InsertAt(index int, value T) error {
	n := node.NewNode(value)
	var err error
	l.withLock(func() {
		if index == 0 {
			n.Next, l.head = l.head, &n
		} else {
			current, next := l.head, l.head
			for index > 0 {
				current, next, index = next, next.Next, index-1
				if next == nil { // Handle index out of bounds
					err = errors.New("index out of bounds")
					return
				}
			}

			current.Next, n.Next = &n, next
		}
	})
	return err
}

func (l *LinkedList[T]) Prepend(value T) {
	n := node.NewNode(value)
	l.withLock(func() {
		n.Next = l.head
		l.head = &n
	})
}

func (l *LinkedList[T]) IsEmpty() bool {
	return l.head == nil
}

func (l *LinkedList[T]) All(yield func(int, T) bool) {
	i := 0
	l.withLock(func() {
		for current := l.head; current != nil; current = current.Next {
			if !yield(i, current.Inner) {
				break
			}
			i++
		}
	})
}

func (l *LinkedList[T]) Size() int {
	i := 0
	l.withLock(func() {
		current := l.head
		if current != nil {
			i++
			for current.Next != nil {
				current = current.Next
				i++
			}
		}
	})

	return i
}

func (l *LinkedList[T]) IndexOf(value T) int {
	if l.head == nil {
		return -1
	}
	index, i := -1, 0
	l.withLock(func() {
		for current := l.head; current != nil; current = current.Next {
			if current.Inner == value {
				index = i
				return
			}
			i++
		}
	})

	return index
}

func (l *LinkedList[T]) GetLast() (T, bool) {
	var (
		v T
		b bool
	)

	if l.head == nil {
		return v, b
	}

	l.withLock(func() {
		current := l.head
		for current.Next != nil {
			current = current.Next
		}
		v = current.Inner
		b = true
	})
	return v, b
}

func (l *LinkedList[T]) GetFirst() (T, bool) {
	var (
		v T
		b bool
	)

	l.withLock(func() {
		if l.head != nil {
			v = l.head.Inner
			b = true
		}
	})

	return v, b
}

func (l *LinkedList[T]) GetAt(index int) (T, error) {
	var (
		v   T
		err error
	)

	l.withLock(func() {
		i := 0
		for current := l.head; current != nil; current = current.Next {
			if i == index {
				v = current.Inner
				return
			}
			i++
		}
		err = errors.New("index out of range")
	})

	return v, err
}

func (l *LinkedList[T]) RemoveLast() (T, bool) {
	var (
		v T
		b bool
	)

	l.withLock(func() {
		prev, current := l.head, l.head
		for current.Next != nil {
			prev = current
			current = current.Next
		}
		prev.Next = nil
		v = current.Inner
		b = true
	})

	return v, b
}

func (l *LinkedList[T]) RemoveFirst() (T, bool) {
	var (
		v T
		b bool
	)
	l.withLock(func() {
		current := l.head
		if current != nil {
			v = current.Inner
			b = true
			l.head = current.Next
		}
	})
	return v, b
}

func (l *LinkedList[T]) Remove(val T) bool {
	b := false

	l.withLock(func() {
		if l.head != nil {

			if l.head.Inner == val {
				l.head = l.head.Next
				b = true
				return
			}

			prev := l.head
			for current := l.head; current != nil; current = current.Next {
				if current.Inner == val {
					b = true
					prev.Next = current.Next
					return
				}
				prev = current
			}
		}
	})
	return b
}

func (l *LinkedList[T]) RemoveAt(idx int) (T, error) {
	var (
		v   T
		err error
		b   bool
	)

	if idx == 0 {
		v, b = l.RemoveFirst()
		if !b {
			err = errors.New("index out of bounds")
		}
		return v, err
	}

	l.withLock(func() {
		i := 0
		prev := l.head
		for current := l.head; current != nil; current = current.Next {
			if i == idx {
				v = current.Inner
				prev.Next = current.Next
				return
			}
			prev = current
			i++
		}
		err = errors.New("index out of bounds")
	})

	return v, err
}

func (l *LinkedList[T]) Clear() {
	l.withLock(func() {
		// this should detach head and trigger GC at some point
		l.head = nil
	})
}

func (l *LinkedList[T]) Items() []T {
	var rs []T
	for _, s := range l.All {
		rs = append(rs, s)
	}
	return rs
}

func (l *LinkedList[T]) AddAll(items ...T) *LinkedList[T] {
	sl := NewLinkedList[T]()
	cur := sl.head
	if cur == nil && len(items) > 0 {
		f := node.NewNode(items[0])
		sl.head = &f
		cur = sl.head
		items = items[1:]
	}

	for _, item := range items {
		n := node.NewNode(item)
		cur.Next = &n
		cur = &n
	}
	l.Join(&sl)
	return l
}

func (l *LinkedList[T]) UnmarshalJSON(data []byte) error {
	var aux []T

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	l.AddAll(aux...)

	return nil
}

func (l *LinkedList[T]) MarshalJSON() ([]byte, error) {
	items := l.Items()
	return json.Marshal(items)
}
