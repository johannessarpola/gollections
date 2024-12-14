package gollections

import (
	"errors"
	"fmt"
	"testing"
)

func TestUnwrap(t *testing.T) {
	var rss []Result[int]

	rss = append(rss, Wrap(1, nil))
	rss = append(rss, Wrap(2, nil))
	rss = append(rss, Wrap(3, nil))
	rss = append(rss, Wrap(0, errors.New("ping pong computer is broke")))

	ecount := 0
	r := UnwrapResults(rss, func(err error) {
		ecount++
	})

	if len(r) == len(rss) {
		t.Errorf("error should have been filtered and not included in result values")
	}

	if ecount == 0 {
		t.Errorf("error should have been filtered and sent into the handler")
	}

}

func TestFanout(t *testing.T) {
	var rss []Result[int]

	rss = append(rss, Wrap(1, nil))
	rss = append(rss, Wrap(2, nil))
	rss = append(rss, Wrap(3, nil))
	rss = append(rss, Wrap(0, errors.New("ping pong computer is broke")))

	r, e := FanOut(rss)

	if len(r) == len(rss) {
		t.Errorf("error should have been filtered and not included in result values")
	}

	if len(e) == 0 {
		t.Errorf("error should have been filtered into err list")
	}

}

func TestFlatMap(t *testing.T) {
	// Success case
	r1 := Wrap(10, nil)
	result := FlatMap(r1, func(v int) Result[string] {
		return Wrap(fmt.Sprintf("Value: %d", v), nil)
	})
	if !result.OK() || result.Val != "Value: 10" {
		t.Errorf("expected Val: 'Value: 10', got Val: %v, Err: %v", result.Val, result.Err)
	}

	// Error propagation case
	r2 := Wrap(0, errors.New("initial error"))
	result = FlatMap(r2, func(v int) Result[string] {
		return Wrap("This should not be executed", nil)
	})
	if result.OK() || result.Err.Error() != "initial error" {
		t.Errorf("expected Err: 'initial error', got Val: %v, Err: %v", result.Val, result.Err)
	}
}

func TestMap(t *testing.T) {
	// Success case
	r1 := Wrap(10, nil)
	result := Map(r1, func(v int) string {
		return fmt.Sprintf("Value: %d", v)
	})
	if !result.OK() || result.Val != "Value: 10" {
		t.Errorf("expected Val: 'Value: 10', got Val: %v, Err: %v", result.Val, result.Err)
	}

	// Error propagation case
	r2 := Wrap(0, errors.New("initial error"))
	result = Map(r2, func(v int) string {
		return "This should not be executed"
	})
	if result.OK() || result.Err.Error() != "initial error" {
		t.Errorf("expected Err: 'initial error', got Val: %v, Err: %v", result.Val, result.Err)
	}
}

func TestMapError(t *testing.T) {
	// Success case (no change)
	r1 := Wrap(10, nil)
	result := MapError(r1, func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})
	if !result.OK() || result.Val != 10 {
		t.Errorf("expected Val: 10, got Val: %v, Err: %v", result.Val, result.Err)
	}

	// Error transformation case
	r2 := Wrap(0, errors.New("initial error"))
	result = MapError(r2, func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})
	if result.OK() || result.Err.Error() != "wrapped: initial error" {
		t.Errorf("expected Err: 'wrapped: initial error', got Val: %v, Err: %v", result.Val, result.Err)
	}
}

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
