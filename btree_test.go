package gollections

import (
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
