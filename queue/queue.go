package gollections

import (
	"sync"

	"github.com/johannessarpola/gollections/internal/node"
)

// Queue FIFO data structure
type Queue[T comparable] struct {
	head *node.Node[T]
	last *node.Node[T]
	mu   sync.Mutex
}

func NewQueue[T comparable]() Queue[T] {
	return Queue[T]{}
}

func (r *Queue[T]) withLock(f func()) {
	defer r.mu.Unlock()
	r.mu.Lock()
	f()
}

func (q *Queue[T]) Enqueue(value T) {
	q.withLock(func() {
		n := node.NewNode(value)
		if q.last != nil {
			q.last.Next = &n
		}
		q.last = &n

		if q.head == nil {
			q.head = &n
		}
	})
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var (
		val T
		ok  bool
	)

	q.withLock(func() {
		if q.head != nil {
			val, _ = q.head.Get()
			ok = true
			q.head = q.head.Next

			if q.head == nil {
				q.last = nil
			}
		}
	})

	return val, ok
}

func (q *Queue[T]) Peek() (T, bool) {
	var (
		val     T
		nonZero bool
	)

	q.withLock(func() {
		if q.head != nil {
			val, nonZero = q.head.Get()
		}
	})
	return val, nonZero
}

func (q *Queue[T]) Size() int {
	i := 0
	q.withLock(func() {
		current := q.head
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

func (q *Queue[T]) IsEmpty() bool {
	b := false
	q.withLock(func() {
		b = q.head == nil
	})
	return b
}
