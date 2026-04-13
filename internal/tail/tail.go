package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

const pollInterval = 200 * time.Millisecond

// Tailer reads new lines appended to a file, similar to `tail -f`.
type Tailer struct {
	path   string
	file   *os.File
	reader *bufio.Reader
}

// New opens the file at path and seeks to the end, ready to tail new content.
func New(path string) (*Tailer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}
	return &Tailer{
		path:   path,
		file:   f,
		reader: bufio.NewReader(f),
	}, nil
}

// Lines returns a channel that emits new lines as they are appended to the file.
// The channel is closed when ctx is cancelled or an unrecoverable error occurs.
func (t *Tailer) Lines(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for {
			line, err := t.reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					return
				}
				// No new data yet — wait before retrying.
				select {
				case <-ctx.Done():
					return
				case <-time.After(pollInterval):
					continue
				}
			}
			if line != "" {
				select {
				case ch <- line:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return ch
}

// Close releases the underlying file handle.
func (t *Tailer) Close() error {
	return t.file.Close()
}
