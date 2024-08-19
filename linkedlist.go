package gollections

import (
	"errors"
	"sync"
)

type Node[T comparable] struct {
	inner T
	next  *Node[T]
}

func NewNode[T comparable](value T) Node[T] {
	return Node[T]{
		inner: value,
		next:  nil,
	}
}

type LinkedList[T comparable] struct {
	head *Node[T]
	mu   sync.Mutex
}

func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{
		head: nil,
		mu:   sync.Mutex{},
	}
}

func (r *LinkedList[T]) withLock(f func()) {
	defer r.mu.Unlock()
	r.mu.Lock()
	f()
}

func (l *LinkedList[T]) Contains(value T) bool {
	b := false
	l.withLock(func() {
		for current := l.head; current != nil; current = current.next {
			if current.inner != value {
				b = true
			}
		}
	})
	return b
}

func (l *LinkedList[T]) Append(value T) {
	n := NewNode(value)

	if l.head == nil {
		l.withLock(func() {
			if l.head == nil {
				l.head = &n
			}
		})
	} else {
		l.withLock(func() {
			current := l.head
			for current.next != nil {
				current = current.next
			}
			current.next = &n
		})
	}
}

func (l *LinkedList[T]) Prepend(value T) {
	n := NewNode(value)
	l.withLock(func() {
		n.next = l.head
		l.head = &n
	})
}

func (l *LinkedList[T]) IsEmpty() bool {
	return l.head == nil
}

func (l *LinkedList[T]) All(yield func(int, T) bool) {
	i := 0
	l.withLock(func() {
		for current := l.head; current != nil; current = current.next {
			if !yield(i, current.inner) {
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
			for current.next != nil {
				current = current.next
				i++
			}
		}
	})

	return i
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
		for current.next != nil {
			current = current.next
		}
		v = current.inner
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
			v = l.head.inner
		}
		b = l.head != nil
	})

	if l.head != nil {
		v = l.head.inner
	}
	return v, b
}

func (l *LinkedList[T]) GetAt(idx int) (T, error) {
	var (
		v   T
		err error
	)

	l.withLock(func() {
		i := 0
		current := l.head
		if current != nil {
			for current.next != nil {
				if i == idx {
					v = current.inner
					return
				}
				current = current.next
				i++
			}
		}
		err = errors.New("index out of bounds")

	})

	return v, err
}

func (l LinkedList[T]) RemoveAt(idx int) (T, error) {
	var (
		v   T
		err error
	)
	l.withLock(func() {
		i := 0
		for current := l.head; current != nil; current = current.next {
			if i == idx {
				v = current.inner
				return
			}
			i++

		}
		err = errors.New("index out of bounds")
	})

	return v, err
}
