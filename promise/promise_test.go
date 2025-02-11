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

func TestPromise_Catch(t *testing.T) {
	tests := []struct {
		name      string
		promise   Promise[string]
		expectErr bool
		expected  string
		errorMsg  string
	}{
		{
			name:      "Successful promise should not trigger Catch",
			promise:   Resolve("Hello, World!"),
			expectErr: false,
			expected:  "Hello, World!",
		},
		{
			name:      "Failed promise should trigger Catch",
			promise:   Reject[string](errors.New("Something went wrong")),
			expectErr: true,
			errorMsg:  "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var caughtError error
			var resultValue string

			// Attach Catch to promise
			done := make(chan struct{})
			tt.promise.Catch(func(err error) {
				caughtError = err
				close(done)
			})

			// Read from the promise
			select {
			case res := <-tt.promise:
				if res.OK() {
					resultValue = res.Value()
				}
				if res.IsErr() {
					caughtError = res.Err()
				}
			case <-time.After(1 * time.Second):
				t.Fatal("Test timed out")
			}

			if tt.expectErr {
				assert.Error(t, caughtError)
				assert.Equal(t, tt.errorMsg, caughtError.Error())
			} else {
				assert.NoError(t, caughtError)
				assert.Equal(t, tt.expected, resultValue)
			}
		})
	}
}

func TestAll(t *testing.T) {
	type testCase struct {
		name       string
		promises   []Promise[int]
		expected   []result.Result[int]
		shouldFail bool
	}

	testCases := map[string]testCase{
		"all-success": {
			promises: []Promise[int]{
				Resolve(1),
				Resolve(2),
				Resolve(3),
			},
			expected: []result.Result[int]{
				result.NewOk(1),
				result.NewOk(2),
				result.NewOk(3),
			},
		},
		"one-error": {
			promises: []Promise[int]{
				Resolve(1),
				Reject[int](errors.New("error")),
				Resolve(3),
			},
			expected: []result.Result[int]{
				result.NewOk(1),
				result.NewErr[int](errors.New("error")),
				result.NewOk(3),
			},
			shouldFail: false,
		},
		"all-error": {
			promises: []Promise[int]{
				Reject[int](errors.New("error")),
				Reject[int](errors.New("another error")),
			},
			expected: []result.Result[int]{
				result.NewErr[int](errors.New("error")),
				result.NewErr[int](errors.New("another error")),
			},
			shouldFail: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rssr := All(context.TODO(), tc.promises...).Wait()
			if rssr.Err() != nil && tc.shouldFail == false {
				t.Errorf("expected no error, got %v", rssr.Err())
			}

			rss := rssr.Value()
			if len(rss) != len(tc.expected) {
				t.Errorf("expected %d results, got %d", len(tc.expected), len(rss))
			}
			for _, rs := range rss {
				found := false
				v, e := rs.Get()
				for _, exp := range tc.expected {
					ev, ee := exp.Get()
					if v == ev {
						found = true
						break
					}
					if ee != nil && e != nil {
						e1 := ee.Error()
						e2 := e.Error()
						if e1 == e2 {
							found = true
							break
						}
					}
				}
				if !found {
					t.Errorf("could not find stuff")
				}
			}
		})
	}
}
