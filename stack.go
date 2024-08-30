package gollections

import (
	"sync"
)

type Stack[T comparable] struct {
	head *Node[T]
	mu   sync.Mutex
}

func NewStack[T comparable]() Stack[T] {
	return Stack[T]{}
}

func (r *Stack[T]) withLock(f func()) {
	defer r.mu.Unlock()
	r.mu.Lock()
	f()
}

func (r *Stack[T]) Peek() (T, bool) {
	var (
		v T
		b bool
	)
	r.withLock(func() {
		if r.head != nil {
			v, b = r.head.Get()
		}
	})

	return v, b
}

func (r *Stack[T]) Push(e T) {
	r.withLock(func() {
		n := NewNode(e)
		n.next = r.head
		r.head = &n
	})
}

func (r *Stack[T]) Pop() (T, bool) {
	var (
		v T
		b bool
	)
	r.withLock(func() {
		if r.head != nil {
			h := r.head
			r.head = h.next
			v, _ = h.Get()
			b = true
		}
	})
	return v, b
}

func (q *Stack[T]) IsEmpty() bool {
	b := false
	q.withLock(func() {
		b = q.head == nil
	})
	return b
}
