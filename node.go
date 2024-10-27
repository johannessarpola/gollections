package gollections

import "fmt"

type Node[T comparable] struct {
	inner T
	next  *Node[T]
	prev  *Node[T]
}

func (n *Node[T]) Get() (T, bool) {
	var zv T
	return n.inner, n.inner != zv
}

func (n *Node[T]) String() string {
	return fmt.Sprintf("%v", n.inner)
}

func (n *Node[T]) HasPrev() bool {
	return n.prev != nil
}

func (n *Node[T]) HasNext() bool {
	return n.next != nil
}

func NewNode[T comparable](value T) Node[T] {
	return Node[T]{
		inner: value,
		next:  nil,
	}
}
