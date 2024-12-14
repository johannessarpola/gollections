package result

// Result is a wrapper for the value and an error.
type Result[T any] struct {
	Val T
	Err error
}

func New[T any](val T, err error) Result[T] {
	return Result[T]{Val: val, Err: err}
}

// OK returns true if the result is ok.
func (r Result[T]) OK() bool {
	return r.Err == nil
}

// Value returns the value of the result.
func (r Result[T]) Value() T {
	return r.Val
}

// Error Standard error interface
func (r Result[T]) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	return ""
}

// OrElse returns the value or a fallback value if the Result is not ok.
func (r Result[T]) OrElse(fallback T) T {
	if r.OK() {
		return r.Val
	}
	return fallback
}

// OrElseFunc returns the value or computes a fallback value if the Result is not ok.
func (r Result[T]) OrElseFunc(fallback func() T) T {
	if r.OK() {
		return r.Val
	}
	return fallback()
}

// Get returns the result as a tuple of value and error.
func (r Result[T]) Get() (T, error) {
	return r.Val, r.Err
}
