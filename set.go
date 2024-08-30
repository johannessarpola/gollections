package gollections

import "sync"

// Set
type Set[T comparable] struct {
	internal map[T]struct{}
	mu       sync.Mutex
}

// NewSet creates a new Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{internal: make(map[T]struct{})}
}

func (s *Set[T]) withLock(f func()) {
	defer s.mu.Unlock()
	s.mu.Lock()
	f()
}

func (s *Set[T]) Add(value T) {
	s.internal[value] = struct{}{}
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
