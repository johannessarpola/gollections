package flagvar

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/johannessarpola/gollections/optional"
)

// OptionalFlag is a generic flag wrapper for optional values
type OptionalFlag[T any] struct {
	opt optional.Optional[T]
}

// Ensure OptionalFlag[T] implements flag.Value interface
var (
	_ flag.Value = (*OptionalFlag[string])(nil)
	_ flag.Value = (*OptionalFlag[int])(nil)
	_ flag.Value = (*OptionalFlag[bool])(nil)
)

/*
 * funny thing about type system is that you need to alias bool heere.
 * Otherwise it works in the BoolFlag() as wotherwise it is handled as a alias for a parameter.
 * So in a sense you need to separate the the two.
 */
type boolean = bool

// IsBoolFlag makes it possible to specify `-enabled` without a value (implies true)
func (f *OptionalFlag[boolean]) IsBoolFlag() bool {
	return true
}

// Set parses the flag value based on T
func (f *OptionalFlag[T]) Set(s string) error {
	switch any(*new(T)).(type) {
	case string:
		f.opt = optional.Some(any(s).(T))

	case int:
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", s)
		}
		f.opt = optional.Some(any(v).(T))

	case bool:
		// If the argument starts with `-`, assume it was passed without a value -> default `true`
		if s == "" || s[0] == '-' {
			f.opt = optional.Some(any(true).(T))
			return nil
		}
		v, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", s)
		}
		f.opt = optional.Some(any(v).(T))

	default:
		return fmt.Errorf("unsupported flag type: %T", *new(T))
	}

	return nil
}

// String returns the stored value as a string
func (f *OptionalFlag[T]) String() string {
	if f.opt.IsPresent() {
		return fmt.Sprintf("%v", f.opt.Get())
	}
	return "unset"
}
