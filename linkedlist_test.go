package gollections

import "testing"

func TestLinkedList1(t *testing.T) {
	l := NewLinkedList[int]()

	l.Append(1)
	l.Append(2)
	l.Append(3)

	l.Prepend(99)

	if v, ok := l.GetLast(); v != 3 || !ok {
		t.Errorf("expected last to be %d and %t but got %d and %t", 3, true, v, ok)
	}

	if l.Size() != 4 {
		t.Errorf("expected size to be %d got %d", 4, l.Size())
	}

	if v, err := l.Get(0); v != 99 || err != nil {
		t.Errorf("expected element at 0 to be %d but got %d with error %v", 99, v, err)
	}

	if _, err := l.Get(10); err == nil {
		t.Error("expected error but got nil")
	}

}
