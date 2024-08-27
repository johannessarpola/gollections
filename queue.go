package gollections

type Queue[T comparable] struct {
	internal LinkedList[T]
}

func NewQueue[T comparable]() Queue[T] {
	return Queue[T]{
		internal: NewLinkedList[T](),
	}
}

func (q *Queue[T]) Enqueue(value T) {
	q.internal.Append(value)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	return q.internal.RemoveLast()
}

func (q *Queue[T]) Peek() (T, bool) {
	return q.internal.GetLast()
}

func (q *Queue[T]) Size() int {
	return q.internal.Size()
}

func (q *Queue[T]) IsEmpty() bool {
	return q.internal.IsEmpty()
}

func (q *Queue[T]) Clear() {
	q.internal.Clear()
}
