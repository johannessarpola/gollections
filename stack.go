package gollections

import (
	"sync"
)

type Stack[T any] struct {
	inner []T
	mu    sync.Mutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{inner: make([]T, 0)}
}

func (r *Stack[T]) mutate(f func()) {
	defer r.mu.Unlock()
	r.mu.Lock()
	f()
}

func (r *Stack[T]) Peek() (T, bool) {
	var (
		rs T
		b  bool
	)
	b = true
	r.mutate(func() {
		if r.inner == nil || len(r.inner) == 0 {
			b = false
			return
		}
		rs = r.inner[len(r.inner)-1]
	})
	return rs, b
}

func (r *Stack[T]) Push(e T) bool {
	b := true
	r.mutate(func() {
		r.inner = append(r.inner, e)
	})
	return b
}

func (r *Stack[T]) Pop() (T, bool) {
	var (
		rs T
		b  bool
	)
	b = true
	r.mutate(func() {
		if r.inner == nil || len(r.inner) == 0 {
			b = false
			return
		}
		l := len(r.inner)
		last, rest := r.inner[l-1], r.inner[:l-1]
		r.inner = rest
		rs = last
	})
	return rs, b
}
