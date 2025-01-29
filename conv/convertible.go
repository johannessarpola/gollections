// conv contains various methods for converting values from one type to another.
package conv

import (
	"fmt"
	"strconv"
)

// Convertible type constraint
type Convertible interface {
	int | bool | string
}

// Convert function
func Convert[T Convertible](value T) string {
	switch v := any(value).(type) {
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	default:
		return ""
	}
}

// Parse function: Converts string to int, bool, or returns the original string
func Parse[T Convertible](input string) (T, error) {
	var zero T

	// Try to parse as int
	if v, err := strconv.Atoi(input); err == nil {
		if any(zero) == int(0) {
			return any(v).(T), nil
		}
	}

	// Try to parse as bool
	if v, err := strconv.ParseBool(input); err == nil {
		if any(zero) == bool(false) {
			return any(v).(T), nil
		}
	}

	// Default: Return as string
	if any(zero) == "" {
		return any(input).(T), nil
	}

	return zero, fmt.Errorf("unsupported type: %T", zero)
}
