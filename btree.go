package gollections

import (
	"cmp"
	"sync"
)

type BinaryTree[T cmp.Ordered] struct {
	head *Node[T]
	mu   sync.Mutex
}

func NewBinaryTree[T cmp.Ordered]() *BinaryTree[T] {
	return &BinaryTree[T]{}
}

func (r *BinaryTree[T]) withLock(f func()) {
	defer r.mu.Unlock()
	r.mu.Lock()
	f()
}

func insert[T cmp.Ordered](root *Node[T], value T) *Node[T] {
	n := NewNode(value)
	if root == nil {
		return &n
	}
	v, _ := root.Get()
	if value < v {
		root.prev = insert(root.prev, value)
	} else {
		root.next = insert(root.next, value)
	}

	return root
}

func (b *BinaryTree[T]) Insert(value T) {
	b.withLock(func() {
		b.head = insert(b.head, value)
	})
}
