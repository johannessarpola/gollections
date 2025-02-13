package optional

// NewExisting wraps value into optional with Exist=true
func NewExisting[T any](value T) Optional[T] {
	return Optional[T]{Value: value, Exist: true}
}

func EmptyInt() Optional[int] {
	return Empty[int]()
}

func EmptyString() Optional[string] {
	return Empty[string]()
}

func Some[T any](val T) Optional[T] {
	return NewExisting(val)
}

func None[T any]() Optional[T] {
	return Empty[T]()
}

// NewString creates a new optional with value
func NewString(val string) Optional[string] {
	return Optional[string]{
		Value: val,
		Exist: val != "",
	}
}
