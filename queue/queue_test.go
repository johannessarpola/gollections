package gollections

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue_Enqueue(t *testing.T) {
	q := Queue[int]{}

	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	if q.head == nil || q.head.Inner != 10 {
		t.Errorf("Expected head value to be 10, but got %v", q.head)
	}

	if q.head.Next == nil || q.head.Next.Inner != 20 {
		t.Errorf("Expected head value to be 20, but got %v", q.head.Next)
	}

	if q.last == nil || q.last.Inner != 30 {
		t.Errorf("Expected last value to be 30, but got %v", q.last)
	}

	if q.last.Next != nil {
		t.Errorf("Expected last values next to be nil, but got %v", q.last.Next)
	}
}

func TestQueue_Dequeue(t *testing.T) {
	q := Queue[int]{}

	q.Enqueue(10)
	if q.last == nil || q.last.Inner != 10 || q.head.Inner != 10 {
		t.Errorf("Expected last value to be 10, but got %v", q.last)
	}

	q.Enqueue(20)
	if q.last == nil || q.last.Inner != 20 || q.head.Inner != 10 {
		t.Errorf("Expected last value to be 20, but got %v", q.last)
	}

	q.Enqueue(30)
	if q.last == nil || q.last.Inner != 30 || q.head.Inner != 10 {
		t.Errorf("Expected last value to be 30, but got %v", q.last)
	}

	v, ok := q.Dequeue()
	if !ok || v != 10 || q.head.Inner == 10 {
		t.Errorf("Expected to dequeue 10, but got %v", v)
	}

	v, ok = q.Dequeue()
	if !ok || v != 20 || q.head.Inner == 20 {
		t.Errorf("Expected to dequeue 20, but got %v", v)
	}

	v, ok = q.Dequeue()
	if !ok || v != 30 {
		t.Errorf("Expected to dequeue 30, but got %v", v)
	}

	v, ok = q.Dequeue()
	if ok || v != 0 {
		t.Errorf("Expected to dequeue from an empty queue and get zero value, but got %v", v)
	}

	if s := q.Size(); s != 0 {
		t.Errorf("Expected to have 0 elements, but got %v", s)
	}

	if q.last != nil {
		t.Errorf("Expected last to be nil, but got %v", q.last)
	}

	if q.head != nil {
		t.Errorf("Expected head to be nil, but got %v", q.head)
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	q := Queue[int]{}

	if !q.IsEmpty() {
		t.Error("Expected queue to be empty initially")
	}

	q.Enqueue(10)
	if q.IsEmpty() {
		t.Error("Expected queue to be non-empty after enqueue")
	}

	q.Dequeue()
	if !q.IsEmpty() {
		t.Error("Expected queue to be empty after dequeueing the only element")
	}
}

func TestQueue_ThreadSafety(t *testing.T) {
	q := Queue[int]{}
	wg := sync.WaitGroup{}
	cnt := 100
	pcntr := atomic.Int32{}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			for i := 0; i < cnt; i++ {
				q.Enqueue(i)
				pcntr.Add(1)
			}
		}()
	}

	dcntr := atomic.Int32{}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			for i := 0; i < cnt; i++ {
				q.Dequeue()
				dcntr.Add(1)
			}
		}()
	}

	// fixme for some reason this is unreliable test
	time.Sleep(1 * time.Second)
	wg.Wait()
	if dcntr.Load() != pcntr.Load() {
		t.Errorf("Expected %v, but got %v", pcntr.Load(), dcntr.Load())
	}
}
