// params package contains handling for optional parameters in httpRequests.
package params

// ParamOpts is a generic struct that holds optional parameters of type T.
type ParamOpts[T any] struct {
	DV T // default value
}

// ParamOpt represents a functional option for configuring the call.
type ParamOpt[T any] func(*ParamOpts[T])

// WithDefault sets the default value.
func WithDefault[T any](value T) ParamOpt[T] {
	return func(o *ParamOpts[T]) {
		o.DV = value
	}
}
