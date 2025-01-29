package params

import (
	"net/url"

	"github.com/johannessarpola/gollections/conv"
)

func argHandler[T any](opts []ParamOpt[T]) ParamOpts[T] {
	args := ParamOpts[T]{}
	for _, opt := range opts {
		opt(&args)
	}
	return args
}

// Check if a parameter exists in the request URL
func HasParam(field string, values url.Values) bool {
	return values.Get(field) != ""
}

// GetParam gets a parameter from the request URL and parses it to the specified type
func GetParam[T conv.Convertible](field string, values url.Values, paramOpts ...ParamOpt[T]) T {
	dv := argHandler(paramOpts).DV
	if p := values.Get(field); p != "" {
		v, err := conv.Parse[T](p)
		if err == nil {
			return dv
		}
		return v
	}
	return dv
}
