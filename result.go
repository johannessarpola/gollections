package gollections

// Result is a wrapper for the value and an error.
type Result[T any] struct {
	Val T
	Err error
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

// Wrap creates a new result from the given value and error.
func Wrap[T any](data T, err error) Result[T] {
	return Result[T]{Val: data, Err: err}
}

// Wrap creates a new result from return values of function f.
func WrapFunc[T any](f func() (T, error)) Result[T] {
	d, e := f()
	return Result[T]{Val: d, Err: e}
}

// Get returns the result as a tuple of value and error.
func (r Result[T]) Get() (T, error) {
	return r.Val, r.Err
}

// UnwrapResults unwraps list of results and calls callback on errored results,
func UnwrapResults[T any](results []Result[T], onError func(err error)) []T {
	var list []T
	for _, res := range results {
		v, err := res.Get()
		if res.OK() {
			list = append(list, v)
		}
		if err != nil {
			onError(err)
		}
	}

	return list
}

// FanOut returns a slice of T and an error slice.
func FanOut[T any](results []Result[T]) ([]T, []error) {
	var (
		list []T
		errs []error
	)
	for _, res := range results {
		v, err := res.Get()
		if res.OK() {
			list = append(list, v)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	return list, errs
}

// Map applies a transformation function to the value of a Result if it is ok.
func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.OK() {
		return Result[U]{Val: f(r.Val), Err: nil}
	}
	// Return a zero value of U with the existing error.
	return Result[U]{Val: *new(U), Err: r.Err}
}

// MapError applies a transformation function to the error if the result is not ok.
func MapError[T any](r Result[T], f func(error) error) Result[T] {
	if !r.OK() {
		return Result[T]{Val: r.Val, Err: f(r.Err)}
	}
	return r
}

// FlatMap applies a function that returns a new Result if the current Result is ok.
func FlatMap[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.OK() {
		return f(r.Val)
	}
	return Result[U]{Val: *new(U), Err: r.Err}
}

// Combine merges a slice of Results into a single Result containing a slice of values if all are ok.
func Combine[T any](results []Result[T]) Result[[]T] {
	var values []T
	for _, res := range results {
		if !res.OK() {
			return Result[[]T]{Err: res.Err}
		}
		values = append(values, res.Val)
	}
	return Result[[]T]{Val: values, Err: nil}
}
