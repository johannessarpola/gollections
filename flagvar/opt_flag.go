package flagvar

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/johannessarpola/gollections/optional"
)

// OptFlag is a generic flag wrapper for optional values
type OptFlag[T any] struct {
	Value optional.Optional[T]
}

// Ensure OptionalFlag[T] implements flag.Value interface
var (
	_ flag.Value = (*OptFlag[string])(nil)
	_ flag.Value = (*OptFlag[int])(nil)
	_ flag.Value = (*OptFlag[bool])(nil)
)

/*
 * funny thing about type system is that you need to alias bool heere.
 * Otherwise it works in the BoolFlag() as wotherwise it is handled as a alias for a parameter.
 */
type boolean = bool

// IsBoolFlag makes it possible to specify `-enabled` without a value (implies true)
func (f *OptFlag[boolean]) IsBoolFlag() bool {
	return true
}

// Set parses the flag value based on T
func (f *OptFlag[T]) Set(s string) error {
	switch any(*new(T)).(type) {
	case string:
		f.Value = optional.Some(any(s).(T))

	case int:
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", s)
		}
		f.Value = optional.Some(any(v).(T))

	case bool:
		// If the argument starts with `-`, assume it was passed without a value -> default `true`
		if s == "" || s[0] == '-' {
			f.Value = optional.Some(any(true).(T))
			return nil
		}
		v, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", s)
		}
		f.Value = optional.Some(any(v).(T))

	default:
		return fmt.Errorf("unsupported flag type: %T", *new(T))
	}

	return nil
}

// String returns the stored value as a string
func (f *OptFlag[T]) String() string {
	if f.Value.IsPresent() {
		return fmt.Sprintf("%v", f.Value.Get())
	}
	return "unset"
}
