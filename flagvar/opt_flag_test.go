package flagvar

import (
	"flag"
	"testing"

	"github.com/johannessarpola/gollections/optional"
)

// TestOptionalFlag tests OptionalFlag[T] with string, int, and bool types.
func TestOptionalFlag(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectStr  optional.Optional[string]
		expectInt  optional.Optional[int]
		expectBool optional.Optional[bool]
	}{
		{
			name:       "All flags provided",
			args:       []string{"-name=Alice", "-age=30", "-enabled=true"},
			expectStr:  optional.Some("Alice"),
			expectInt:  optional.Some(30),
			expectBool: optional.Some(true),
		},
		{
			name:       "Only string and int flags",
			args:       []string{"-name=Bob", "-age=42"},
			expectStr:  optional.Some("Bob"),
			expectInt:  optional.Some(42),
			expectBool: optional.None[bool](),
		},
		{
			name:       "Only boolean flag (implicit true)",
			args:       []string{"-enabled"},
			expectStr:  optional.None[string](),
			expectInt:  optional.None[int](),
			expectBool: optional.Some(true), // Implicitly true
		},
		{
			name:       "Boolean flag explicitly set to false",
			args:       []string{"-enabled=false"},
			expectStr:  optional.None[string](),
			expectInt:  optional.None[int](),
			expectBool: optional.Some(false),
		},
		{
			name:       "Boolean flag explicitly set to true",
			args:       []string{"-enabled=true"},
			expectStr:  optional.None[string](),
			expectInt:  optional.None[int](),
			expectBool: optional.Some(true),
		},
		{
			name:       "Only int flag",
			args:       []string{"-age=99"},
			expectStr:  optional.None[string](),
			expectInt:  optional.Some(99),
			expectBool: optional.None[bool](),
		},
		{
			name:       "No flags provided",
			args:       []string{},
			expectStr:  optional.None[string](),
			expectInt:  optional.None[int](),
			expectBool: optional.None[bool](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag set for each test case
			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			var strFlag OptFlag[string]
			var intFlag OptFlag[int]
			var boolFlag OptFlag[bool]

			fs.Var(&strFlag, "name", "A string flag")
			fs.Var(&intFlag, "age", "An integer flag")
			fs.Var(&boolFlag, "enabled", "A boolean flag (true if present)")

			// Parse arguments
			err := fs.Parse(tt.args)
			if err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Validate string flag
			if strFlag.Value != tt.expectStr {
				t.Errorf("expected string flag = %v, got %v", tt.expectStr, strFlag.Value)
			}

			// Validate int flag
			if intFlag.Value != tt.expectInt {
				t.Errorf("expected int flag = %v, got %v", tt.expectInt, intFlag.Value)
			}

			// Validate bool flag
			if boolFlag.Value != tt.expectBool {
				t.Errorf("expected bool flag = %v, got %v", tt.expectBool, boolFlag.Value)
			}
		})
	}
}
