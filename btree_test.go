package gollections

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	bt := NewBinaryTree[int]()

	bt.Insert(99)

	if v, b := bt.head.Get(); v != 99 || !b {
		t.Errorf("head should be %v, want %v", 99, v)
	}

	bt.Insert(77)
	if v, b := bt.head.prev.Get(); v != 77 || !b {
		t.Errorf("head left should be %d got %v", 77, v)
	}

	bt.Insert(101)
	if v, b := bt.head.next.Get(); v != 101 || !b {
		t.Errorf("head right should be %d got %v", 101, v)
	}

	bt.Insert(1)
	if v, b := bt.head.prev.prev.Get(); v != 1 || !b {
		t.Errorf("head left.left should be %d got %v", 1, v)
	}

}

type TravelsalOrder = string

const (
	InOrder    TravelsalOrder = "inOrder"
	PreOrder                  = "preOrder"
	PostOrder                 = "postOrder"
	LevelOrder                = "levelOrder"
)

func (bt *BinaryTree[T]) resolveTravelsalFunc(travelsalOrder TravelsalOrder) func(yield func(int, T) bool) {
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

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{name: "size-1", input: []int{1, 2, 3, 4, 5}, expected: 5},
		{name: "size-2", input: []int{1}, expected: 1},
		{name: "size-3", input: []int{1, 2, 3}, expected: 3},
		{name: "size-4", input: []int{}, expected: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}

			if bt.Size() != test.expected {
				t.Errorf("got %v, want %v", bt.Size(), test.expected)
			}

		})
	}
}

func TestInorder(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
		to       TravelsalOrder
	}{
		{name: "preOrder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{99, 77, 33, 90, 101}, to: PreOrder},
		{name: "preOrder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, to: PreOrder},
		{name: "preOrder-3", input: []int{1, 99, 101, 2}, expected: []int{1, 99, 2, 101}, to: PreOrder},
		{name: "preOrder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{5, 3, 2, 4, 7, 6, 8}, to: PreOrder},
		{name: "preOrder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{50, 30, 20, 40, 70, 60, 80}, to: PreOrder},
		{name: "preOrder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{15, 10, 8, 12, 20, 17, 25}, to: PreOrder},
		{name: "preOrder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{10, 5, 1, 7, 40, 30, 50}, to: PreOrder},
		{name: "preOrder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{12, 8, 5, 10, 15, 13, 18}, to: PreOrder},

		{name: "inOrder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{33, 77, 90, 99, 101}, to: InOrder},
		{name: "inOrder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, to: InOrder},
		{name: "inOrder-3", input: []int{1, 99, 101, 2}, expected: []int{1, 2, 99, 101}, to: InOrder},
		{name: "inOrder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{2, 3, 4, 5, 6, 7, 8}, to: InOrder},
		{name: "inOrder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{20, 30, 40, 50, 60, 70, 80}, to: InOrder},
		{name: "inOrder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{8, 10, 12, 15, 17, 20, 25}, to: InOrder},
		{name: "inOrder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{1, 5, 7, 10, 30, 40, 50}, to: InOrder},
		{name: "inOrder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{5, 8, 10, 12, 13, 15, 18}, to: InOrder},

		{name: "postOrder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{33, 90, 77, 101, 99}, to: PostOrder},
		{name: "postOrder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{5, 4, 3, 2, 1}, to: PostOrder},
		{name: "postOrder-3", input: []int{1, 99, 101, 2}, expected: []int{2, 101, 99, 1}, to: PostOrder},
		{name: "postOrder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{2, 4, 3, 6, 8, 7, 5}, to: PostOrder},
		{name: "postOrder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{20, 40, 30, 60, 80, 70, 50}, to: PostOrder},
		{name: "postOrder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{8, 12, 10, 17, 25, 20, 15}, to: PostOrder},
		{name: "postOrder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{1, 7, 5, 30, 50, 40, 10}, to: PostOrder},
		{name: "postOrder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{5, 10, 8, 13, 18, 15, 12}, to: PostOrder},

		{name: "levelorder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{99, 77, 101, 33, 90}, to: LevelOrder},
		{name: "levelorder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, to: LevelOrder},
		{name: "levelorder-3", input: []int{1, 99, 101, 1}, expected: []int{1, 99, 1, 101}, to: LevelOrder},
		{name: "levelorder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{5, 3, 7, 2, 4, 6, 8}, to: LevelOrder},
		{name: "levelorder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{50, 30, 70, 20, 40, 60, 80}, to: LevelOrder},
		{name: "levelorder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{15, 10, 20, 8, 12, 17, 25}, to: LevelOrder},
		{name: "levelorder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{10, 5, 40, 1, 7, 30, 50}, to: LevelOrder},
		{name: "levelorder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{12, 8, 15, 5, 10, 13, 18}, to: LevelOrder},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}

			fmt.Println("tree representation: ")
			fmt.Print(bt.String())

			var rs []int
			travelsalFunc := bt.resolveTravelsalFunc(test.to)
			for _, v := range travelsalFunc {
				rs = append(rs, v)
			}

			fmt.Printf("%s travelsal was:\n%v\n", test.to, rs)

			if !reflect.DeepEqual(rs, test.expected) {
				t.Errorf("got %v, want %v", rs, test.expected)
			}

		})
	}
}
