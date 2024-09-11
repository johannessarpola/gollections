package gollections

import (
	"cmp"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

type BinaryTree[T cmp.Ordered] struct {
	head *Node[T]
	mu   sync.Mutex
}

func NewBinaryTree[T cmp.Ordered]() BinaryTree[T] {
	return BinaryTree[T]{}
}

func (bt *BinaryTree[T]) withLock(f func()) {
	defer bt.mu.Unlock()
	bt.mu.Lock()
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

func (bt *BinaryTree[T]) Insert(values ...T) {
	bt.withLock(func() {
		for _, v := range values {
			bt.head = insert(bt.head, v)
		}
	})
}

func preOrder[T cmp.Ordered](root *Node[T], i *atomic.Int32, yield func(int, T) bool) {
	if root == nil {
		return
	}

	// Yield the current node
	index := int(i.Add(1) - 1)
	if !yield(index, root.inner) {
		return
	}

	// Traverse the left subtree
	preOrder(root.prev, i, yield)

	// Traverse the right subtree
	preOrder(root.next, i, yield)

}

func inOrder[T cmp.Ordered](root *Node[T], i *atomic.Int32, yield func(int, T) bool) {
	if root == nil {
		return
	}

	inOrder(root.prev, i, yield)

	// Yield the current node
	index := int(i.Add(1) - 1)
	if !yield(index, root.inner) {
		return
	}

	inOrder(root.next, i, yield)
}

func postOrder[T cmp.Ordered](root *Node[T], i *atomic.Int32, yield func(int, T) bool) {
	if root == nil {
		return
	}

	postOrder(root.prev, i, yield)

	postOrder(root.next, i, yield)
	// Yield the current node
	index := int(i.Add(1) - 1)
	if !yield(index, root.inner) {
		return
	}

}

func levelOrder[T cmp.Ordered](root *Node[T], yield func(int, T) bool) {
	if root == nil {
		return
	}

	queue := make([]*Node[T], 0)
	queue = append(queue, root)

	i := 0

	for len(queue) > 0 {
		// dequeue the first currentNode
		currentNode := queue[0] // deque element
		queue = queue[1:]       // remove first element

		// yield the currentNode value
		if !yield(i, currentNode.inner) {
			return
		}
		i++

		// add left child to the queue
		if currentNode.prev != nil {
			queue = append(queue, currentNode.prev)
		}

		// add right child to the queue
		if currentNode.next != nil {
			queue = append(queue, currentNode.next)
		}
	}
}

// treeHeight computes the height of the binary tree.
func treeHeight[T comparable](root *Node[T]) int {
	if root == nil {
		return 0
	}
	leftHeight := treeHeight(root.prev)
	rightHeight := treeHeight(root.next)

	// The height of the tree is the maximum of the two subtrees + 1 (for the current root).
	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// Postorder left-root-right
func (bt *BinaryTree[T]) Postorder(yield func(int, T) bool) {
	bt.withLock(func() {
		postOrder(bt.head, &atomic.Int32{}, yield)
	})
}

// InOrder left-root-right
func (bt *BinaryTree[T]) InOrder(yield func(int, T) bool) {
	bt.withLock(func() {
		inOrder(bt.head, &atomic.Int32{}, yield)
	})
}

// PreOrder root-left-right
func (bt *BinaryTree[T]) PreOrder(yield func(int, T) bool) {
	bt.withLock(func() {
		preOrder(bt.head, &atomic.Int32{}, yield)
	})
}

// LeverOrder breadth first
func (bt *BinaryTree[T]) LeverOrder(yield func(int, T) bool) {
	bt.withLock(func() {
		levelOrder(bt.head, yield)
	})
}

func (bt *BinaryTree[T]) Height() int {
	return treeHeight(bt.head)
}

// String returns a string visualization of the binary tree.
func (bt *BinaryTree[T]) String() string {
	var sb strings.Builder
	bt.withLock(func() {
		if bt.head == nil {
			sb.WriteString("empty")
			return
		}

		bt.visualizeNode(bt.head, "", true, &sb)
	})

	return sb.String()
}

// visualizeNode helps in the recursive visualization of the binary tree.
func (bt *BinaryTree[T]) visualizeNode(node *Node[T], prefix string, isTail bool, sb *strings.Builder) {
	if node == nil {
		return
	}

	// Append current node value
	sb.WriteString(prefix)
	if isTail {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("├── ")
	}
	sb.WriteString(fmt.Sprintf("%v\n", node.inner))

	// Prepare the prefix for child nodes
	childPrefix := prefix
	if isTail {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Recurse on children nodes
	hasPrev := node.prev != nil
	hasNext := node.next != nil

	if hasPrev {
		bt.visualizeNode(node.prev, childPrefix, !hasNext, sb)
	}

	if hasNext {
		bt.visualizeNode(node.next, childPrefix, true, sb)
	}
}
