package gollections

import (
	"encoding/json"
	"slices"
	"sync"

	"github.com/johannessarpola/gollections/internal/node"
)

type Stack[T comparable] struct {
	head *node.Node[T]
	mu   sync.Mutex
}

func NewStack[T comparable]() Stack[T] {
	return Stack[T]{}
}

func (s *Stack[T]) withLock(f func()) {
	defer s.mu.Unlock()
	s.mu.Lock()
	f()
}

func (s *Stack[T]) Peek() (T, bool) {
	var (
		v T
		b bool
	)
	s.withLock(func() {
		if s.head != nil {
			v, b = s.head.Get()
		}
	})

	return v, b
}

func (s *Stack[T]) Push(e T) {
	s.withLock(func() {
		n := node.NewNode(e)
		n.Next = s.head
		s.head = &n
	})
}

func (s *Stack[T]) Pop() (T, bool) {
	var (
		v T
		b bool
	)
	s.withLock(func() {
		if s.head != nil {
			h := s.head
			s.head = h.Next
			v, _ = h.Get()
			b = true
		}
	})
	return v, b
}

func (s *Stack[T]) IsEmpty() bool {
	b := false
	s.withLock(func() {
		b = s.head == nil
	})
	return b
}

func (s *Stack[T]) PushAll(items ...T) {
	if len(items) == 0 {
		return
	}

	first := node.NewNode(items[0])
	iteratee := &first
	rest := items[1:]

	for _, item := range rest {
		node := node.NewNode(item)
		node.Next = iteratee
		iteratee = &node
	}

	s.withLock(func() {
		if s.head == nil {
			s.head = iteratee
			return
		}
		iteratee.Next = s.head
		s.head = iteratee
	})
}

func (s *Stack[T]) PopAll() []T {
	var (
		rs       []T
		iteratee T
	)

	ok := true
	for ok {
		iteratee, ok = s.Pop()
		if ok {
			rs = append(rs, iteratee)
		}
	}

	return rs
}

func (s *Stack[T]) UnmarshalJSON(data []byte) error {
	var aux []T

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.PushAll(aux...)

	return nil
}

func (s *Stack[T]) MarshalJSON() ([]byte, error) {
	itemsReversed := s.PopAll()
	slices.Reverse(itemsReversed)
	return json.Marshal(itemsReversed)
}
