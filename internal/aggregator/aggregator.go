// Package aggregator combines log entries from multiple readers,
// applies filter criteria, and writes matching entries to an output.
package aggregator

import (
	"sync"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/reader"
)

// Source pairs a named service reader with its log stream.
type Source struct {
	Name   string
	Reader *reader.Reader
}

// Aggregator fans-in log entries from multiple sources, filters them,
// and forwards matching entries to a single output writer.
type Aggregator struct {
	sources  []Source
	criteria filter.Criteria
	out      *output.Writer
}

// New creates an Aggregator for the given sources, filter criteria, and output writer.
func New(sources []Source, criteria filter.Criteria, out *output.Writer) *Aggregator {
	return &Aggregator{
		sources:  sources,
		criteria: criteria,
		out:      out,
	}
}

// Run starts reading from all sources concurrently and writes matching
// log entries to the output. It blocks until all sources are exhausted.
func (a *Aggregator) Run() error {
	type result struct {
		entry reader.LogEntry
		err   error
	}

	entriesCh := make(chan result, len(a.sources)*8)

	var wg sync.WaitGroup
	for _, src := range a.sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()
			for {
				entry, err := s.Reader.Next()
				if err != nil {
					// nil entry signals EOF or unrecoverable error
					entriesCh <- result{err: err}
					return
				}
				entriesCh <- result{entry: entry}
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(entriesCh)
	}()

	for res := range entriesCh {
		if res.err != nil {
			continue // skip EOF / parse errors from individual sources
		}
		if a.criteria.Match(res.entry) {
			if err := a.out.Write(res.entry); err != nil {
				return err
			}
		}
	}
	return nil
}
