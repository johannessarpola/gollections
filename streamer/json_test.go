package streamer

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
)

// Helper function to read and unmarshal JSON lines
func readJSONLines[T any](reader io.Reader) ([]T, error) {
	var results []T
	decoder := json.NewDecoder(reader)

	for {
		var obj T
		if err := decoder.Decode(&obj); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		results = append(results, obj)
	}

	return results, nil
}

type testStruct struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func TestArrayStreamReader(t *testing.T) {
	tests := []struct {
		name      string
		data      []any
		expectErr bool
		expected  []any
	}{
		{
			name: "Valid data",
			data: []any{
				testStruct{ID: 1, Name: "Alice"},
				testStruct{ID: 2, Name: "Bob", Tags: []string{"admin", "user"}},
				testStruct{ID: 3, Name: "Charlie", Tags: []string{"user"}},
			},
			expectErr: false,
			expected: []any{
				testStruct{ID: 1, Name: "Alice"},
				testStruct{ID: 2, Name: "Bob", Tags: []string{"admin", "user"}},
				testStruct{ID: 3, Name: "Charlie", Tags: []string{"user"}},
			},
		},
		{
			name:      "Empty data",
			data:      []any{},
			expectErr: false,
			expected:  []any{},
		},
		{
			name: "Invalid data type",
			data: []any{
				make(chan int),
			},
			expectErr: true,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewArrayStreamReader(tt.data)
			onErrorCalled := false
			reader.Start(t.Context(), func(err error) {
				if !tt.expectErr {
					t.Fatalf("Unexpected error: %v", err)
				}
				onErrorCalled = true
			})

			var buf bytes.Buffer
			_, err := io.Copy(&buf, reader)

			if tt.expectErr && !onErrorCalled {
				t.Fatalf("Expected error, but none occurred")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectErr {
				results, err := readJSONLines[testStruct](&buf)
				if err != nil {
					t.Fatalf("Failed to read JSON lines: %v", err)
				}

				if len(results) != len(tt.expected) {
					t.Errorf("Expected %d JSON objects, got %d", len(tt.expected), len(results))
				}

				for i, result := range results {
					expected := tt.expected[i]
					if !reflect.DeepEqual(result, expected) {
						t.Errorf("Expected %v, got %v", expected, result)
					}
				}
			}
		})
	}
}
