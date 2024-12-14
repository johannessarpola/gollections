package result

// Wrap creates a new result from the given value and error.
func Wrap[T any](data T, err error) Result[T] {
	return Result[T]{Val: data, Err: err}
}

// Wrap creates a new result from return values of function f.
func WrapFunc[T any](f func() (T, error)) Result[T] {
	d, e := f()
	return Result[T]{Val: d, Err: e}
}

// UnwrapSlice unwraps list of results and calls callback on errored results,
func UnwrapSlice[T any](results []Result[T], onError func(err error)) []T {
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
