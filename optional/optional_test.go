package optional

import (
	"encoding/json"
	"testing"
)

func TestNewOptional(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		wantExist bool
	}{
		{
			name:      "create optional with value",
			value:     10,
			wantExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewExisting(tt.value)
			if opt.Value != tt.value {
				t.Errorf("NewOptional() got = %v, want %v", opt.Value, tt.value)
			}
			if opt.Exist != tt.wantExist {
				t.Errorf("NewOptional() gotExist = %v, want %v", opt.Exist, tt.wantExist)
			}
		})
	}
}

func TestEmptyOptional(t *testing.T) {
	tests := []struct {
		name      string
		wantExist bool
	}{
		{
			name:      "create empty optional",
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Empty[int]()
			if opt.Exist != tt.wantExist {
				t.Errorf("EmptyOptional() gotExist = %v, want %v", opt.Exist, tt.wantExist)
			}
		})
	}
}

func TestIsPresent(t *testing.T) {
	tests := []struct {
		name string
		opt  Optional[int]
		want bool
	}{
		{
			name: "value is present",
			opt:  NewExisting(5),
			want: true,
		},
		{
			name: "value is not present",
			opt:  Empty[int](),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsPresent(); got != tt.want {
				t.Errorf("IsPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		opt       Optional[int]
		want      int
		wantPanic bool
	}{
		{
			name:      "value exists",
			opt:       NewExisting(10),
			want:      10,
			wantPanic: false,
		},
		{
			name:      "value does not exist",
			opt:       Empty[int](),
			wantPanic: true, // This will cause a panic.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("Get() panicked unexpectedly")
					}
				}
			}()

			got := tt.opt.Get()

			if !tt.wantPanic && got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		opt          Optional[int]
		defaultValue int
		want         int
	}{
		{
			name:         "value exists, return value",
			opt:          NewExisting(10),
			defaultValue: 0,
			want:         10,
		},
		{
			name:         "value does not exist, return default",
			opt:          Empty[int](),
			defaultValue: 42,
			want:         42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.GetOrDefault(tt.defaultValue); got != tt.want {
				t.Errorf("GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionalMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Optional[any]
		expected string
	}{
		{
			"String value",
			New[any]("hello", true),
			`"hello"`,
		},
		{
			"Int value",
			New[any](42, true),
			`42`,
		},
		{
			"Bool value",
			New[any](true, true),
			`true`,
		},
		{
			"Zero int but Exist true",
			New[any](0, true),
			`0`,
		},
		{
			"Empty string but Exist true",
			New[any]("", false),
			"null",
		},
		{
			"Exist false (should be omitted)",
			New[any]("", false),
			"null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, data)
			}
		})
	}
}

func TestOptionalUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Optional[any]
	}{
		{"String present", `"Hello"`, Optional[any]{Value: "Hello", Exist: true}},
		{"Int present", `42`, Optional[any]{Value: float64(42), Exist: true}}, // JSON numbers decode to float64 by default
		{"Bool present", `true`, Optional[any]{Value: true, Exist: true}},
		{"Zero int", `0`, Optional[any]{Value: float64(0), Exist: true}},
		{"Empty string", `""`, Optional[any]{Value: "", Exist: true}},
		{"Null (missing field equivalent)", `null`, Optional[any]{Value: nil, Exist: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Optional[any]
			err := json.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if result.Exist != tt.expected.Exist || result.Value != tt.expected.Value {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

type Parent struct {
	Name   Optional[string]
	Age    Optional[int]
	Height Optional[float64]
}

func TestParentStructUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Parent
	}{
		{
			"All fields present",
			`{"name":"Alice","age":30}`,
			Parent{
				Name:   Optional[string]{Value: "Alice", Exist: true},
				Age:    Optional[int]{Value: 30, Exist: true},
				Height: Optional[float64]{Exist: false},
			},
		},
		{
			"Missing age field",
			`{"name":"Alice"}`,
			Parent{
				Name:   Optional[string]{Value: "Alice", Exist: true},
				Age:    Optional[int]{Exist: false},
				Height: Optional[float64]{Exist: false},
			},
		},
		{
			"Empty JSON (all fields missing)",
			`{}`,
			Parent{
				Name:   Optional[string]{Exist: false},
				Age:    Optional[int]{Exist: false},
				Height: Optional[float64]{Exist: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Parent
			err := json.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if result.Name != tt.expected.Name || result.Age != tt.expected.Age {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
