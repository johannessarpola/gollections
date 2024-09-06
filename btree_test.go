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

type TravelsalOrder = int

const (
	Inorder TravelsalOrder = iota
	Preorder
	Postorder
)

func (bt *BinaryTree[T]) resolveTravelsalFunc(travelsalOrder TravelsalOrder) func(yield func(int, T) bool) {
	switch travelsalOrder {
	case Inorder:
		return bt.Inorder
	case Preorder:
		return bt.Preorder
	case Postorder:
		return bt.Postorder
	default:
		panic(fmt.Sprintf("unknown travelsal order %v", travelsalOrder))
	}
}

func TestInorder(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
		to       TravelsalOrder
	}{
		{name: "preorder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{99, 77, 33, 90, 101}, to: Preorder},
		{name: "preorder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, to: Preorder},
		{name: "preorder-3", input: []int{1, 99, 101, 2}, expected: []int{1, 99, 2, 101}, to: Preorder},
		{name: "preorder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{5, 3, 2, 4, 7, 6, 8}, to: Preorder},
		{name: "preorder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{50, 30, 20, 40, 70, 60, 80}, to: Preorder},
		{name: "preorder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{15, 10, 8, 12, 20, 17, 25}, to: Preorder},
		{name: "preorder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{10, 5, 1, 7, 40, 30, 50}, to: Preorder},
		{name: "preorder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{12, 8, 5, 10, 15, 13, 18}, to: Preorder},

		{name: "inorder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{33, 77, 90, 99, 101}, to: Inorder},
		{name: "inorder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, to: Inorder},
		{name: "inorder-3", input: []int{1, 99, 101, 2}, expected: []int{1, 2, 99, 101}, to: Inorder},
		{name: "inorder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{2, 3, 4, 5, 6, 7, 8}, to: Inorder},
		{name: "inorder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{20, 30, 40, 50, 60, 70, 80}, to: Inorder},
		{name: "inorder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{8, 10, 12, 15, 17, 20, 25}, to: Inorder},
		{name: "inorder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{1, 5, 7, 10, 30, 40, 50}, to: Inorder},
		{name: "inorder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{5, 8, 10, 12, 13, 15, 18}, to: Inorder},

		{name: "postorder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{33, 90, 77, 101, 99}, to: Postorder},
		{name: "postorder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{5, 4, 3, 2, 1}, to: Postorder},
		{name: "postorder-3", input: []int{1, 99, 101, 2}, expected: []int{2, 101, 99, 1}, to: Postorder},
		{name: "postorder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{2, 4, 3, 6, 8, 7, 5}, to: Postorder},
		{name: "postorder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{20, 40, 30, 60, 80, 70, 50}, to: Postorder},
		{name: "postorder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{8, 12, 10, 17, 25, 20, 15}, to: Postorder},
		{name: "postorder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{1, 7, 5, 30, 50, 40, 10}, to: Postorder},
		{name: "postorder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{5, 10, 8, 13, 18, 15, 12}, to: Postorder},
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

			fmt.Printf("resolveTravelsalFunc order was:\n%v\n", rs)

			if !reflect.DeepEqual(rs, test.expected) {
				t.Errorf("got %v, want %v", rs, test.expected)
			}

		})
	}
}
