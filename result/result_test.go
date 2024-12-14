package result

import (
	"errors"
	"testing"
)

func TestOrElse(t *testing.T) {
	// Success case
	r1 := Wrap(10, nil)
	val := r1.OrElse(42)
	if val != 10 {
		t.Errorf("expected 10, got %v", val)
	}

	// Error fallback case
	r2 := Wrap(0, errors.New("initial error"))
	val = r2.OrElse(42)
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestOrElseFunc(t *testing.T) {
	// Success case
	r1 := Wrap(10, nil)
	val := r1.OrElseFunc(func() int {
		return 42
	})
	if val != 10 {
		t.Errorf("expected 10, got %v", val)
	}

	// Error fallback case
	r2 := Wrap(0, errors.New("initial error"))
	val = r2.OrElseFunc(func() int {
		return 42
	})
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}
