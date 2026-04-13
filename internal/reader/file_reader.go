package reader

import (
	"context"

	"github.com/yourorg/logpipe/internal/tail"
)

// FileReader wraps a Tailer to produce LogEntry values from a tailed file.
type FileReader struct {
	service string
	tailer  *tail.Tailer
}

// NewFileReader creates a FileReader that tails the file at path, attributing
// entries to the given service name.
func NewFileReader(path, service string) (*FileReader, error) {
	t, err := tail.New(path)
	if err != nil {
		return nil, err
	}
	return &FileReader{service: service, tailer: t}, nil
}

// Next returns a channel of LogEntry values parsed from new lines in the file.
// The channel is closed when ctx is cancelled.
func (fr *FileReader) Next(ctx context.Context) <-chan LogEntry {
	out := make(chan LogEntry)
	go func() {
		defer close(out)
		for line := range fr.tailer.Lines(ctx) {
			entry := parseLine(line, fr.service)
			select {
			case out <- entry:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Close releases the underlying file handle.
func (fr *FileReader) Close() error {
	return fr.tailer.Close()
}
