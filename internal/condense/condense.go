package condense

import (
	"errors"
	"fmt"
	"time"

	"github.com/your-org/logpipe/internal/reader"
)

// Condenser merges consecutive log entries from the same service whose
// messages share a common prefix into a single representative entry.
// This is useful for reducing noise from repeated partial-line flushes.

// Options configures the Condenser.
type Options struct {
	// MinPrefix is the minimum shared prefix length required to condense.
	MinPrefix int
	// MaxAge is the maximum time an open group is held before it is flushed.
	MaxAge time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MinPrefix: 8,
		MaxAge:    2 * time.Second,
	}
}

// Condenser groups and merges qualifying log entries.
type Condenser struct {
	opts    Options
	pending map[string]*group
	clock   func() time.Time
}

type group struct {
	base   reader.LogEntry
	count  int
	opened time.Time
	prefix string
}

// New returns a Condenser with the given options.
func New(opts Options) (*Condenser, error) {
	if opts.MinPrefix < 1 {
		return nil, errors.New("condense: MinPrefix must be at least 1")
	}
	if opts.MaxAge <= 0 {
		return nil, errors.New("condense: MaxAge must be positive")
	}
	return &Condenser{
		opts:    opts,
		pending: make(map[string]*group),
		clock:   time.Now,
	}, nil
}

// Apply attempts to merge entry into an open group for its service.
// It returns the flushed entry when a group is closed, or nil when the
// entry has been absorbed. Callers must also call Flush periodically.
func (c *Condenser) Apply(entry reader.LogEntry) *reader.LogEntry {
	now := c.clock()
	key := entry.Service

	if g, ok := c.pending[key]; ok {
		shared := commonPrefix(g.prefix, entry.Message)
		if len(shared) >= c.opts.MinPrefix && now.Sub(g.opened) < c.opts.MaxAge {
			g.count++
			g.prefix = shared
			return nil
		}
		out := c.closeGroup(g)
		c.pending[key] = c.newGroup(entry, now)
		return &out
	}

	c.pending[key] = c.newGroup(entry, now)
	return nil
}

// Flush closes all open groups older than MaxAge and returns their entries.
func (c *Condenser) Flush() []reader.LogEntry {
	now := c.clock()
	var out []reader.LogEntry
	for key, g := range c.pending {
		if now.Sub(g.opened) >= c.opts.MaxAge {
			out = append(out, c.closeGroup(g))
			delete(c.pending, key)
		}
	}
	return out
}

func (c *Condenser) newGroup(entry reader.LogEntry, now time.Time) *group {
	return &group{base: entry, count: 1, opened: now, prefix: entry.Message}
}

func (c *Condenser) closeGroup(g *group) reader.LogEntry {
	out := g.base
	if g.count > 1 {
		out.Message = fmt.Sprintf("%s … (%dx)", g.prefix, g.count)
	}
	return out
}

func commonPrefix(a, b string) string {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:n]
}
