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

	c := l.Contains(99)
	if c == true {
		t.Errorf("expected removed element to not be contained but got %t", c)
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

	c2 := l.Contains(-99)
	if c2 == true {
		t.Errorf("expected element to be not contained but got %t", c2)
	}

	err = l.InsertAt(0, -66)
	if err != nil {
		t.Error("expected no error but got ", err)
	}

	if v, ok := l.GetFirst(); v != -66 || !ok {
		t.Errorf("expected first element to be %d and %t but got %d and %t", -66, true, v, ok)
	}

	c3 := l.Contains(-66)
	if c3 != true {
		t.Errorf("expected element to be contained but got %t", c3)
	}

	if v, ok := l.GetFirst(); v != -66 || !ok {
		t.Errorf("expected first to be %d but got %d and %t", -66, v, ok)
	}

	err = l.InsertAt(3, 55)
	if err != nil {
		t.Error("expected no error but got ", err)
	}

	if v, err := l.GetAt(3); v != 55 || err != nil {
		t.Errorf("expected %d to be %d but got %d and %s", 3, 55, v, err)
	}

	err = l.InsertAt(99, 33)
	if err == nil {
		t.Error("expected error but got nil")
	}

}
