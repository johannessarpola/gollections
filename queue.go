package gollections

import "sync"

// Queue FIFO data structure
type Queue[T comparable] struct {
	head *Node[T]
	last *Node[T]
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
		n := NewNode(value)
		if q.last != nil {
			q.last.next = &n
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
			val, ok = q.head.Get()
			q.head = q.head.next

			if q.head == nil {
				q.last = nil
			}
		}

		return

	})

	return val, ok
}

func (q *Queue[T]) Peek() (T, bool) {
	var (
		val T
		ok  bool
	)

	q.withLock(func() {
		if q.head != nil {
			val, ok = q.head.Get()
		}
	})
	return val, ok
}

func (q *Queue[T]) Size() int {
	i := 0
	q.withLock(func() {
		current := q.head
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

func (q *Queue[T]) IsEmpty() bool {
	b := false
	q.withLock(func() {
		b = q.head == nil
		return
	})
	return b
}
