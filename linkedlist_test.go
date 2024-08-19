package gollections

import (
	"testing"
)

func TestLinkedList1(t *testing.T) {
	l := NewLinkedList[int]()

	if !l.IsEmpty() {
		t.Errorf("list is not empty")
	}

	l.Append(1)
	l.Append(2)
	l.Append(77)

	l.Prepend(99)

	if l.IsEmpty() {
		t.Errorf("list is empty")
	}

	if v, ok := l.GetLast(); v != 77 || !ok {
		t.Errorf("expected last to be %d and %t but got %d and %t", 3, true, v, ok)
	}

	if v, ok := l.GetFirst(); v != 99 || !ok {
		t.Errorf("expected first to be %d and %t but got %d and %t", 1, true, v, ok)
	}

	if l.Size() != 4 {
		t.Errorf("expected size to be %d got %d", 4, l.Size())
	}

	if v, err := l.GetAt(0); v != 99 || err != nil {
		t.Errorf("expected element at 0 to be %d but got %d with error %v", 99, v, err)
	}

	if _, err := l.GetAt(10); err == nil {
		t.Error("expected error but got nil")
	}

	e, err := l.RemoveAt(0)
	if err != nil || e == 0 {
		t.Errorf("could not remove at idx %d", 0)
	}
	e2, err := l.RemoveAt(10)
	if err == nil || e2 != 0 {
		t.Errorf("expected error but got %s", err)
	}

	count := 0
	size := l.Size()
	for _, _ = range l.All {
		count++
	}

	if count != size {
		t.Errorf("expected %d elements but got %d", size, count)
	}

	c := l.Contains(99)
	if c != true {
		t.Errorf("expected element to be contained but got %t", c)
	}

}
