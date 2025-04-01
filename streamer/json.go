package streamer

import (
	"context"
	"encoding/json"
	"io"
)

// ArrayStreamReader implements io.Reader for streaming JSON data
type ArrayStreamReader struct {
	data    []any
	encoder *json.Encoder
	reader  *io.PipeReader
	writer  *io.PipeWriter
}

// NewArrayStreamReader initializes a streaming JSON reader for []any
func NewArrayStreamReader(data []any) *ArrayStreamReader {
	pr, pw := io.Pipe()
	encoder := json.NewEncoder(pw)

	jr := &ArrayStreamReader{
		data:    data,
		encoder: encoder,
		reader:  pr,
		writer:  pw,
	}

	return jr
}

func (jr *ArrayStreamReader) Start(ctx context.Context, onError func(error)) {
	go jr.streamJSON(ctx, onError)
}

// Stream JSON data into the writer
func (jr *ArrayStreamReader) streamJSON(ctx context.Context, onError func(error)) {
	defer jr.writer.Close()
	for _, item := range jr.data {
		select {
		case <-ctx.Done():
			return
		default:
			if err := jr.encoder.Encode(item); err != nil {
				onError(err)
				return
			}
		}
	}
}

// Read implements io.Reader
func (jr *ArrayStreamReader) Read(p []byte) (int, error) {
	return jr.reader.Read(p)
}
