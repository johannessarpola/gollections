package comps

import "testing"

func TestUnorderedEquals(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		arr1 []T
		arr2 []T
		want bool
	}
	tests := []testCase[int]{
		{
			name: "unorderdEquals-1",
			arr1: []int{1, 2, 3},
			arr2: []int{3, 2, 1},
			want: true,
		},
		{
			name: "unorderdEquals-2",
			arr1: []int{1, 1, 1},
			arr2: []int{2, 2, 2},
			want: false,
		},
		{
			name: "unorderdEquals-3",
			arr1: []int{1},
			arr2: []int{1},
			want: true,
		},
		{
			name: "unorderdEquals-4",
			arr1: []int{1},
			arr2: []int{1, 1},
			want: false,
		},
		{
			name: "unorderdEquals-5",
			arr1: []int{},
			arr2: []int{},
			want: true,
		},
		{
			name: "unorderdEquals-6",
			arr1: []int{1},
			arr2: []int{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnorderedEquals(tt.arr1, tt.arr2); got != tt.want {
				t.Errorf("UnorderedEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
