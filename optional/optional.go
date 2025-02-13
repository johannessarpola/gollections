package optional

import (
	"encoding/json"
)

// Optional represents a value that may or may not exist
type Optional[T any] struct {
	Value T
	Exist bool
}

// New creates a new Optional with a value
func New[T any](value T, exist bool) Optional[T] {
	return Optional[T]{Value: value, Exist: exist}
}

// Empty creates an empty Optional
func Empty[T any]() Optional[T] {
	return Optional[T]{Exist: false}
}

// IsPresent returns true if the value exists
func (o Optional[T]) IsPresent() bool {
	return o.Exist
}

// IfPresent param func is called if value exists.
func (o Optional[T]) IfPresent(f func(v T)) {
	if o.IsPresent() {
		f(o.Get())
	}
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

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Exist {
		return []byte("null"), nil // Omitting the field entirely
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		o.Exist = true // Field was missing, set Exist = true
		var zero T
		o.Value = zero
		return nil
	}

	o.Exist = true
	return json.Unmarshal(data, &o.Value)
}
