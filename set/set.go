package gollections

import (
	"encoding/json"
	"sync"
)

// Set
type Set[T comparable] struct {
	internal map[T]struct{}
	mu       sync.Mutex
}

// New creates a new Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{internal: make(map[T]struct{})}
}

func (s *Set[T]) withLock(f func()) {
	defer s.mu.Unlock()
	s.mu.Lock()
	f()
}

func (s *Set[T]) AddAll(values ...T) {
	s.withLock(func() {
		for _, v := range values {
			s.internal[v] = struct{}{}
		}
	})
}

func (s *Set[T]) Unset(value T) {
	s.withLock(func() { delete(s.internal, value) })
}

func (s *Set[T]) Add(value T) {
	s.withLock(func() {
		s.internal[value] = struct{}{}
	})
}

func (s *Set[T]) Contains(value T) bool {
	_, exists := s.internal[value]
	return exists
}

func (s *Set[T]) Remove(value T) bool {
	b := s.Contains(value)
	if b {
		delete(s.internal, value)
	}
	return b
}

func (s *Set[T]) Size() int {
	return len(s.internal)
}

func (s *Set[T]) All(yield func(int, T) bool) {
	i := 0
	s.withLock(func() {
		for key := range s.internal {
			if !yield(i, key) {
				break
			}
			i++
		}
	})
}

func (s *Set[T]) Clear() {
	s.internal = make(map[T]struct{})
}

func (s *Set[T]) Items() []T {
	keys := make([]T, 0, len(s.internal))
	for k := range s.internal {
		keys = append(keys, k)
	}
	return keys
}

func (c *Set[T]) UnmarshalJSON(data []byte) error {
	if c.internal == nil {
		c.internal = make(map[T]struct{})
	}

	var aux []T

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.AddAll(aux...)

	return nil
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	items := s.Items()
	return json.Marshal(items)
}
