package promise

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/johannessarpola/gollections/result"
	"github.com/stretchr/testify/assert"
)

func TestThen(t *testing.T) {
	type testCase struct {
		name      string
		input     int
		to        time.Duration
		transform func(context.Context, int) result.Result[int]
		want      int
		wantErr   bool
	}

	testCases := []testCase{
		{
			name:  "test-then-1",
			input: 10,
			transform: func(_ context.Context, i int) result.Result[int] {
				return result.NewOk(i * 10)
			},
			want:    100,
			wantErr: false,
			to:      time.Second * 1,
		},
		{
			name:  "test-then-2",
			input: 10,
			transform: func(_ context.Context, i int) result.Result[int] {
				return result.NewErr[int](errors.New("piip piip"))
			},
			wantErr: true,
			to:      time.Second * 1,
		},
		{
			name:  "test-then-3",
			input: 10,
			to:    10 * time.Millisecond,
			transform: func(ctx context.Context, i int) result.Result[int] {
				select {
				case <-ctx.Done():
					return result.NewErr[int](errors.New("timeout"))
				case <-time.After(1 * time.Second):
					return result.NewOk(i * 2)
				}
			},
			want:    -1,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.to)
			defer cancel()
			start := time.Now()
			p := New[int]().Resolve(ctx, tc.input).Then(ctx, tc.transform)
			rs := p.Wait()
			if tc.wantErr {
				assert.Error(t, rs.Err())
			} else {
				assert.Less(t, time.Since(start), tc.to)
				assert.NoError(t, rs.Err())
				assert.Equal(t, tc.want, rs.Value())
			}
		})
	}
}
