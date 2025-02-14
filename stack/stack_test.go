package gollections

import (
	"encoding/json"
	"reflect"
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

func TestStack_PushAll(t *testing.T) {
	stack := NewStack[int]()
	stack.PushAll(1, 2, 3)
	if v, ok := stack.Pop(); v != 3 || !ok {
		t.Errorf("expected popped value to be 1, got %v", v)
	}

	if v, ok := stack.Pop(); v != 2 || !ok {
		t.Errorf("expected popped value to be 2, got %v", v)
	}

	if v, ok := stack.Pop(); v != 1 || !ok {
		t.Errorf("expected popped value to be 1, got %v", v)
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
	for i := 0; i < 100; i++ {
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

func TestStack_PopAll(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		in   []T
		want []T
	}
	tests := []testCase[int]{
		{name: "popAll-1", in: []int{1, 2, 3}, want: []int{3, 2, 1}},
		{name: "popAll-2", in: []int{}, want: []int{}},
		{name: "popAll-3", in: []int{1}, want: []int{1}},
		{name: "popAll-3", in: []int{6, 7, 8, 9, 10}, want: []int{10, 9, 8, 7, 6}},
	}
	for _, tt := range tests {
		stack := NewStack[int]()
		stack.PushAll(tt.in...)
		t.Run(tt.name, func(t *testing.T) {
			got := stack.PopAll()
			if !reflect.DeepEqual(got, tt.want) && len(tt.want) > 0 {
				t.Errorf("PopAll() = %v, want %v", got, tt.want)
			}

			if len(tt.want) == 0 && len(got) != 0 {
				t.Errorf("PopAll() = expected result to be empty but got %d", len(got))
			}
		})
	}
}

func TestStack_JsonMarshalling(t *testing.T) {
	jsonStr := `[1, 2, 3]`
	var stack Stack[int]
	err := json.Unmarshal([]byte(jsonStr), &stack)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
	}

	popAll := stack.PopAll()
	want := []int{3, 2, 1}
	if !reflect.DeepEqual(popAll, want) {
		t.Errorf("PopAll() = %v, want %v", popAll, want)
	}

	type contained struct {
		Field string        `json:"field"`
		Stack Stack[string] `json:"stack"`
	}

	want2 := []string{"bca", "cba", "abc"}
	var jsonStr2 = `
{
	"field" : "value",
	"stack" : ["abc", "cba", "bca"]
}`

	var c contained
	err = json.Unmarshal([]byte(jsonStr2), &c)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
	}

	if c.Field != "value" {
		t.Errorf("c.Field = %v, want %v", c.Field, "value")
	}

	popAll2 := c.Stack.PopAll()
	if !reflect.DeepEqual(popAll2, want2) {
		t.Errorf("PopAll() = %v, want %v", c.Stack.PopAll(), want2)
	}
}
