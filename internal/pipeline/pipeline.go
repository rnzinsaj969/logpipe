package pipeline

import (
	"context"

	"github.com/user/logpipe/internal/aggregator"
	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/output"
)

// Pipeline wires together an aggregator, a filter criteria, and an output
// writer into a single processing loop.
type Pipeline struct {
	agg      *aggregator.Aggregator
	criteria filter.Criteria
	out      *output.Writer
}

// New creates a Pipeline from the provided components.
func New(agg *aggregator.Aggregator, c filter.Criteria, w *output.Writer) *Pipeline {
	return &Pipeline{
		agg:      agg,
		criteria: c,
		out:      w,
	}
}

// Run reads log entries from the aggregator until the context is cancelled or
// the aggregator is exhausted. Entries that match the criteria are written to
// the output writer. The first write error terminates the loop.
func (p *Pipeline) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		entry, ok := p.agg.Next()
		if !ok {
			return nil
		}

		if !p.criteria.Match(entry) {
			continue
		}

		if err := p.out.Write(entry); err != nil {
			return err
		}
	}
}
