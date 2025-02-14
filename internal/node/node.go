package node

import "fmt"

type Node[T comparable] struct {
	Inner T
	Next  *Node[T]
	Prev  *Node[T]
}

func (n *Node[T]) Get() (T, bool) {
	var zv T
	return n.Inner, n.Inner != zv
}

func (n *Node[T]) String() string {
	return fmt.Sprintf("%v", n.Inner)
}

func (n *Node[T]) HasPrev() bool {
	return n.Prev != nil
}

func (n *Node[T]) HasNext() bool {
	return n.Next != nil
}

func NewNode[T comparable](value T) Node[T] {
	return Node[T]{
		Inner: value,
		Next:  nil,
	}
}
