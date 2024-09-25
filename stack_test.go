package gollections

import (
	"sync"
	"sync/atomic"
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

	// Testing RemoveFirst
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

	ai := atomic.Int32{}
	ai.Store(0)
	// Testing concurrent RemoveFirst
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, ok := stack.Pop()

			if !ok {
				t.Errorf("expected pop to succeed, got %v", ok)
			}
			ai.Add(1)
		}()
	}

	wg.Wait()

	if ai.Load() != 1000 {
		t.Errorf("expected col length to be 1000, got %v", ai.Load())
	}
	if _, ok := stack.Peek(); ok {
		t.Errorf("expected stack to be empty, but got %v", ok)
	}
}

func TestStack_ConcurrentPeek(t *testing.T) {
	stack := NewStack[int]()
	var wg sync.WaitGroup

	// Fill the stack
	for i := 0; i < 100; i++ {
		stack.Push(i)
	}

	last, ok := stack.Peek()
	if !ok {
		t.Errorf("expected successful peek, got %v", ok)
	}

	// Testing concurrent Peek
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, ok := stack.Peek()
			if v != last || !ok {
				t.Errorf("expected successful peek, got %v and %d", ok, v)
			}
		}()
	}

	wg.Wait()

	// Peek should still return the last item pushed
	if v, ok := stack.Peek(); v != 99 || !ok {
		t.Errorf("expected top of stack to be 99, got %v", v)
	}
}
