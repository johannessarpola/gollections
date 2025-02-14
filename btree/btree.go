package btree

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/johannessarpola/gollections/internal/node"
)

type BinaryTree[T cmp.Ordered] struct {
	head           *node.Node[T]
	mu             sync.Mutex
	traversalOrder TraversalOrder // used for json
}

const DefaultTraversal TraversalOrder = PreOrder

type TraversalOrder = string

const (
	InOrder    TraversalOrder = "inOrder"
	PreOrder   TraversalOrder = "preOrder"
	PostOrder  TraversalOrder = "postOrder"
	LevelOrder TraversalOrder = "levelOrder"
)

func (bt *BinaryTree[T]) resolveTravelsalFunc(travelsalOrder TraversalOrder) func(yield func(int, T) bool) {
	switch travelsalOrder {
	case InOrder:
		return bt.InOrder
	case PreOrder:
		return bt.PreOrder
	case PostOrder:
		return bt.Postorder
	case LevelOrder:
		return bt.LeverOrder
	default:
		panic(fmt.Sprintf("unknown travelsal order %v", travelsalOrder))
	}
}

func NewBinaryTree[T cmp.Ordered]() BinaryTree[T] {
	return BinaryTree[T]{}
}

func NewBinaryTreeWithOrder[T cmp.Ordered](order TraversalOrder) BinaryTree[T] {
	return BinaryTree[T]{
		traversalOrder: order,
	}
}

func NewBinaryTreeWithType[T cmp.Ordered](typeOf T, order TraversalOrder) BinaryTree[T] {
	return BinaryTree[T]{
		traversalOrder: order,
	}
}

func (bt *BinaryTree[T]) withLock(f func()) {
	defer bt.mu.Unlock()
	bt.mu.Lock()
	f()
}

func insert[T cmp.Ordered](root *node.Node[T], value T) *node.Node[T] {
	n := node.NewNode(value)
	if root == nil {
		return &n
	}
	v, _ := root.Get()
	if value < v {
		root.Prev = insert(root.Prev, value)
	} else {
		root.Next = insert(root.Next, value)
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

func preOrder[T cmp.Ordered](root *node.Node[T], i *atomic.Int32, yield func(int, T) bool) {
	if root == nil {
		return
	}

	// Yield the current node
	index := int(i.Add(1) - 1)
	if !yield(index, root.Inner) {
		return
	}

	// Traverse the left subtree
	preOrder(root.Prev, i, yield)

	// Traverse the right subtree
	preOrder(root.Next, i, yield)
}

func inOrder[T cmp.Ordered](root *node.Node[T], i *atomic.Int32, yield func(int, T) bool) {
	if root == nil {
		return
	}

	inOrder(root.Prev, i, yield)

	// Yield the current node
	index := int(i.Add(1) - 1)
	if !yield(index, root.Inner) {
		return
	}

	inOrder(root.Next, i, yield)
}

func postOrder[T cmp.Ordered](root *node.Node[T], i *atomic.Int32, yield func(int, T) bool) {
	if root == nil {
		return
	}

	postOrder(root.Prev, i, yield)

	postOrder(root.Next, i, yield)
	// Yield the current node
	index := int(i.Add(1) - 1)
	if !yield(index, root.Inner) {
		return
	}
}

func levelOrder[T cmp.Ordered](root *node.Node[T], yield func(int, T) bool) {
	if root == nil {
		return
	}

	queue := make([]*node.Node[T], 0)
	queue = append(queue, root)

	i := 0

	for len(queue) > 0 {
		// dequeue the first currentNode
		currentNode := queue[0] // deque element
		queue = queue[1:]       // remove first element

		// yield the currentNode value
		if !yield(i, currentNode.Inner) {
			return
		}
		i++

		// add left child to the queue
		if currentNode.Prev != nil {
			queue = append(queue, currentNode.Prev)
		}

		// add right child to the queue
		if currentNode.Next != nil {
			queue = append(queue, currentNode.Next)
		}
	}
}

// treeHeight computes the height of the binary tree.
func treeHeight[T comparable](root *node.Node[T]) int {
	if root == nil {
		return 0
	}
	leftHeight := treeHeight(root.Prev)
	rightHeight := treeHeight(root.Next)

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

// balanceTree builds a balanced binary tree from the sorted slice
func balanceTree[T comparable](values []T, start, end int) *node.Node[T] {
	if start > end {
		return nil
	}

	// middle element as the root
	mid := (start + end) / 2
	node := &node.Node[T]{
		Inner: values[mid],
	}

	// recursively build the left and right subtrees
	node.Prev = balanceTree[T](values, start, mid-1)
	node.Next = balanceTree[T](values, mid+1, end)

	return node
}

func (bt *BinaryTree[T]) Balance() {
	var slice []T
	for _, v := range bt.InOrder {
		slice = append(slice, v)
	}
	slices.Sort(slice)

	bt.withLock(func() {
		bt.head = balanceTree(slice, 0, len(slice)-1)
	})
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

func find[T comparable](root *node.Node[T], predicate func(T, T) bool) (T, bool) {
	var zv T
	if root == nil {
		return zv, false
	}

	best := root.Inner
	lv, lf := find(root.Prev, predicate)

	if predicate(lv, best) && lf {
		best = lv
	}

	rv, rf := find(root.Next, predicate)
	if predicate(rv, best) && rf {
		best = rv
	}
	return best, true
}

func (bt *BinaryTree[T]) FindMax() (T, bool) {
	var (
		rs T
		b  bool
	)
	bt.withLock(func() {
		rs, b = find(bt.head, func(nv T, ov T) bool {
			return nv > ov
		})
	})
	return rs, b
}

func (bt *BinaryTree[T]) FindMin() (T, bool) {
	var (
		rs T
		b  bool
	)
	bt.withLock(func() {
		rs, b = find(bt.head, func(nv T, ov T) bool {
			return nv < ov
		})
	})
	return rs, b
}

func (bt *BinaryTree[T]) Find(predicate func(T, T) bool) (T, bool) {
	var (
		rs T
		b  bool
	)
	bt.withLock(func() {
		rs, b = find(bt.head, predicate)
	})
	return rs, b
}

func search[T cmp.Ordered](root *node.Node[T], target T) (T, bool) {
	var zv T

	if root == nil {
		return zv, false
	}

	if root.Inner == target {
		return root.Inner, true
	}

	if target < root.Inner {
		lrs, lb := search(root.Prev, target)
		if lb {
			return lrs, true
		}
	}

	if target > root.Inner {
		rrs, rb := search(root.Next, target)
		if rb {
			return rrs, true
		}
	}

	return zv, false
}

func (bt *BinaryTree[T]) Search(element T) (T, bool) {
	var (
		rs T
		b  bool
	)
	bt.withLock(func() {
		rs, b = search(bt.head, element)
	})

	return rs, b
}

// visualizeNode helps in the recursive visualization of the binary tree.
func (bt *BinaryTree[T]) visualizeNode(node *node.Node[T], prefix string, isTail bool, sb *strings.Builder) {
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
	sb.WriteString(fmt.Sprintf("%v\n", node.Inner))

	// Prepare the prefix for child nodes
	childPrefix := prefix
	if isTail {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Recurse on children nodes
	hasPrev := node.Prev != nil
	hasNext := node.Next != nil

	if hasPrev {
		bt.visualizeNode(node.Prev, childPrefix, !hasNext, sb)
	}

	if hasNext {
		bt.visualizeNode(node.Next, childPrefix, true, sb)
	}
}

func (bt *BinaryTree[T]) Items() []T {
	f := bt.resolveTravelsalFunc(bt.TraversalOrder())

	var items []T
	for _, v := range f {
		items = append(items, v)
	}

	return items
}

type BinaryTreeJson[T cmp.Ordered] struct {
	Data           []T            `json:"data"`
	TraversalOrder TraversalOrder `json:"traversal_order"`
}

func (bt *BinaryTree[T]) TraversalOrder() TraversalOrder {
	to := bt.traversalOrder
	if bt.traversalOrder == "" {
		to = DefaultTraversal
	}
	return to
}

func (bt *BinaryTree[T]) UnmarshalJSON(data []byte) error {
	var aux BinaryTreeJson[T]

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bt.traversalOrder = aux.TraversalOrder
	bt.Insert(aux.Data...)

	return nil
}

func (bt *BinaryTree[T]) MarshalJSON() ([]byte, error) {
	items := bt.Items()

	aux := BinaryTreeJson[T]{
		Data:           items,
		TraversalOrder: bt.TraversalOrder(),
	}
	return json.Marshal(aux)
}
