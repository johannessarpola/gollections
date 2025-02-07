package result

// Result is a wrapper for the value and an error.
type Result[T any] struct {
	val T
	err error
}

func New[T any](val T, err error) Result[T] {
	return Result[T]{val: val, err: err}
}

func NewOk[T any](val T) Result[T] {
	return Result[T]{val: val, err: nil}
}

func NewErr[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// OK returns true if the result is ok.
func (r Result[T]) OK() bool {
	return r.Err() == nil
}

// Value returns the value of the result.
func (r Result[T]) Value() T {
	return r.val
}

// Error Standard error interface
func (r Result[T]) Err() error {
	return r.err
}

// Error Standard error interface
func (r Result[T]) Error() string {
	if r.Err() != nil {
		return r.Err().Error()
	}
	return ""
}

// Calls `f` if value is present
func (r Result[T]) IfPresent(f func(T) Result[T]) Result[T] {
	if r.OK() {
		return f(r.val)
	}
	return r
}

// OrElse returns the value or a fallback value if the Result is not ok.
func (r Result[T]) OrElse(fallback T) T {
	if r.OK() {
		return r.val
	}
	return fallback
}

// OrElseFunc returns the value or computes a fallback value if the Result is not ok.
func (r Result[T]) OrElseFunc(fallback func() T) T {
	if r.OK() {
		return r.Value()
	}
	return fallback()
}

// Get returns the result as a tuple of value and error.
func (r Result[T]) Get() (T, error) {
	return r.Value(), r.Err()
}
