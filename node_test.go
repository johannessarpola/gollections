package gollections

import (
	"testing"

	"github.com/johannessarpola/gollections/optional"
)

func TestNode(t *testing.T) {
	type testCase[T comparable] struct {
		name    string
		in      T
		prev    optional.Optional[T]
		next    optional.Optional[T]
		want    T
		wantStr string
	}
	intTests := []testCase[int]{
		{name: "node-1", in: 1, want: 1, wantStr: "1"},
		{name: "node-2", in: -13, want: -13, wantStr: "-13"},
		{name: "node-3", in: 3, want: 3, wantStr: "3", next: optional.NewExisting(4)},
		{name: "node-4", in: 4, want: 4, wantStr: "4", prev: optional.NewExisting(3)},
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNode(tt.in)

			if tt.next.IsPresent() {
				nxt := NewNode(tt.next.Get())
				n.next = &nxt
			}

			if tt.prev.IsPresent() {
				prv := NewNode(tt.prev.Get())
				n.prev = &prv
			}

			if n.String() != tt.wantStr {
				t.Errorf("got %v, want %v", n.String(), tt.wantStr)
			}

			if i, b := n.Get(); i != tt.want && !b {
				t.Errorf("got %v, want %v", i, tt.want)
			}

			if (tt.next.IsPresent() && !n.HasNext()) || (!tt.next.IsPresent() && n.HasNext()) {
				t.Errorf("HasNext() returned invalid value")
			}

			if (tt.prev.IsPresent() && !n.HasPrev()) || (!tt.prev.IsPresent() && n.HasPrev()) {
				t.Errorf("HasPrev() returned invalid value")
			}

		})
	}
}
