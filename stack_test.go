package gollections

import (
	"testing"
)

func TestStack_PushAndPop(t *testing.T) {
	stack := NewStack[int]()

	// Testing Push
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if v, ok := stack.Peek(); v != 3 || !ok {
		t.Errorf("expected top of stack to be 3, got %v", v)
	}

	// Testing Pop
	if v, ok := stack.Pop(); v != 3 || !ok {
		t.Errorf("expected popped value to be 3, got %v", v)
	}

	if v, ok := stack.Pop(); v != 2 || !ok {
		t.Errorf("expected popped value to be 2, got %v", v)
	}

	if v, ok := stack.Pop(); v != 1 || !ok {
		t.Errorf("expected popped value to be 1, got %v", v)
	}

	if _, ok := stack.Pop(); ok {
		t.Errorf("expected popped to return not ok, got %v", ok)
	}
}

/*
func TestStack_ConcurrentPushAndPop(t *testing.T) {
	stack := NewStack[int]()
	var wg sync.WaitGroup

	// Testing concurrent Push
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			stack.Push(i)
		}(i)
	}

	wg.Wait()

	// Testing concurrent Pop
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stack.Pop()
		}()
	}

	wg.Wait()

	if stack.Peek() != nil {
		t.Errorf("expected stack to be empty, but got %v", stack.Peek())
	}
}
*/
/*
func TestStack_ConcurrentPeek(t *testing.T) {
	stack := NewStack[int]()
	var wg sync.WaitGroup

	// Fill the stack
	for i := 0; i < 100; i++ {
		stack.Push(i)
	}

	// Testing concurrent Peek
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = stack.Peek()
		}()
	}

	wg.Wait()

	// Peek should still return the last item pushed
	if stack.Peek() != 99 {
		t.Errorf("expected top of stack to be 99, got %v", stack.Peek())
	}
}
*/
