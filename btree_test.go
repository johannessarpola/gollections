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

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{name: "size-1", input: []int{1, 2, 3, 4, 5}, expected: 5},
		{name: "size-2", input: []int{1}, expected: 1},
		{name: "size-3", input: []int{}, expected: 0},
		{name: "size-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: 3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}

			if bt.Height() != test.expected {
				t.Errorf("got %v, want %v", bt.Height(), test.expected)
			}

		})
	}
}

func TestInorder(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
		to       TraversalOrder
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

		{name: "levelOrder-1", input: []int{99, 77, 33, 101, 90}, expected: []int{99, 77, 101, 33, 90}, to: LevelOrder},
		{name: "levelOrder-2", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, to: LevelOrder},
		{name: "levelOrder-3", input: []int{1, 99, 101, 1}, expected: []int{1, 99, 1, 101}, to: LevelOrder},
		{name: "levelOrder-4", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: []int{5, 3, 7, 2, 4, 6, 8}, to: LevelOrder},
		{name: "levelOrder-5", input: []int{50, 30, 20, 40, 70, 60, 80}, expected: []int{50, 30, 70, 20, 40, 60, 80}, to: LevelOrder},
		{name: "levelOrder-6", input: []int{15, 10, 20, 8, 12, 17, 25}, expected: []int{15, 10, 20, 8, 12, 17, 25}, to: LevelOrder},
		{name: "levelOrder-7", input: []int{10, 5, 1, 7, 40, 50, 30}, expected: []int{10, 5, 40, 1, 7, 30, 50}, to: LevelOrder},
		{name: "levelOrder-8", input: []int{12, 8, 15, 5, 10, 13, 18}, expected: []int{12, 8, 15, 5, 10, 13, 18}, to: LevelOrder},
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

func max(i int, i2 int) bool {
	return i > i2
}
func min(i int, i2 int) bool {
	return i < i2
}

func TestFind(t *testing.T) {
	type comparison struct {
		predicate func(int, int) bool
		expected  int
	}
	tests := []struct {
		name  string
		input []int
		comp  []comparison
	}{
		{name: "find-1", input: []int{99, 77, 33, 101, 90}, comp: []comparison{
			{predicate: max, expected: 101},
			{predicate: min, expected: 33},
		}},
		{name: "find-2", input: []int{1, 2, 3, 4, 5}, comp: []comparison{
			{predicate: max, expected: 5},
			{predicate: min, expected: 1},
		}},
		{name: "find-3", input: []int{5, 3, 7, 2, 4, 6, 8}, comp: []comparison{
			{predicate: max, expected: 8},
			{predicate: min, expected: 2},
		}},
		{name: "find-4", input: []int{12, 8, 15, 5, 10, 13, 18}, comp: []comparison{
			{predicate: max, expected: 18},
			{predicate: min, expected: 5},
		}},
		{name: "find-5", input: []int{}, comp: []comparison{
			{predicate: max, expected: 0},
			{predicate: min, expected: 0},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}

			for _, comp := range test.comp {
				rs, b := bt.Find(comp.predicate)
				if rs != comp.expected {
					t.Errorf("got %v, want %v", rs, comp.expected)
				}

				if rs == 0 && b {
					t.Errorf("bool should be false when not found")
				}
			}

		})
	}
}

func TestBinaryTree_FindMax(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{name: "findMax-1", input: []int{99, 77, 33, 101, 90}, expected: 101},
		{name: "findMax-2", input: []int{1, 2, 3, 4, 5}, expected: 5},
		{name: "findMax-3", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: 8},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}
			rs, b := bt.FindMax()
			if rs != test.expected || !b {
				t.Errorf("got %v, want %v", rs, test.expected)
			}
		})
	}
}

func TestBinaryTree_FindMin(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{name: "findMin-1", input: []int{99, 77, 33, 101, 90}, expected: 33},
		{name: "findMin-2", input: []int{1, 2, 3, 4, 5}, expected: 1},
		{name: "findMin-3", input: []int{5, 3, 7, 2, 4, 6, 8}, expected: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}
			rs, b := bt.FindMin()
			if rs != test.expected || !b {
				t.Errorf("got %v, want %v", rs, test.expected)
			}
		})
	}
}

func TestBinaryTree_Search(t *testing.T) {
	tests := []struct {
		name         string
		input        []int
		searchFor    int
		expectedVal  int
		expectedBool bool
	}{
		{name: "search-1", input: []int{99, 77, 33, 101, 90}, searchFor: 33, expectedVal: 33, expectedBool: true},
		{name: "Search-2", input: []int{1, 2, 3, 4, 5}, searchFor: 5, expectedVal: 5, expectedBool: true},
		{name: "search-3", input: []int{5, 3, 7, 2, 4, 6, 8}, searchFor: 8, expectedVal: 8, expectedBool: true},
		{name: "search-4", input: []int{5, 3, 7, 2, 4, 6, 8}, searchFor: 99, expectedVal: 0, expectedBool: false},
		{name: "search-5", input: []int{}, searchFor: 99, expectedVal: 0, expectedBool: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}
			rs, b := bt.Search(test.searchFor)

			if rs != test.expectedVal || b != test.expectedBool {
				t.Errorf("searched for %v, got %v (%v), want %v (%v)", test.searchFor, rs, b, test.expectedVal, test.expectedBool)
			}
		})
	}
}

func TestBinaryTree_Balance(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "balance-1", input: []int{10, 5, 20, 15, 30}, expected: []int{15, 5, 20, 10, 30}},
		{name: "balance-2", input: []int{1, 2, 3, 4, 5}, expected: []int{3, 1, 4, 2, 5}},                           // perfectly balanced
		{name: "balance-3", input: []int{7, 6, 5, 4, 3, 2, 1}, expected: []int{4, 2, 6, 1, 3, 5, 7}},               // balanced after skewed input
		{name: "balance-4", input: []int{30, 20, 40, 10, 25, 35, 50}, expected: []int{30, 20, 40, 10, 25, 35, 50}}, // already balanced
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt := NewBinaryTree[int]()

			for _, v := range test.input {
				bt.Insert(v)
			}

			fmt.Println("tree representation before balancing: ")
			fmt.Print(bt.String())

			bt.Balance()
			fmt.Println("tree representation after balancing: ")
			fmt.Print(bt.String())

			var rs []int
			for _, v := range bt.LeverOrder {
				rs = append(rs, v)
			}

			if !reflect.DeepEqual(rs, test.expected) {
				t.Errorf("got %v, want %v", rs, test.expected)
			}

		})
	}
}
