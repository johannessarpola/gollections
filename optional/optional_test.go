package optional

import (
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
