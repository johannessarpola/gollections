package gollections

// Optional represents a value that may or may not exist
type Optional[T any] struct {
	Value T
	Exist bool
}

// NewOptional creates a new Optional with a value
func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{Value: value, Exist: true}
}

// EmptyOptional creates an empty Optional
func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{Exist: false}
}

// IsPresent returns true if the value exists
func (o Optional[T]) IsPresent() bool {
	return o.Exist
}

// Get returns the value if it exists, otherwise it panics
func (o Optional[T]) Get() T {
	if !o.Exist {
		panic("optional value does not exist")
	}
	return o.Value
}

// GetOrDefault returns the value if it exists, otherwise returns the default value provided
func (o Optional[T]) GetOrDefault(defaultValue T) T {
	if o.Exist {
		return o.Value
	}
	return defaultValue
}
