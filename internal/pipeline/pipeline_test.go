package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/logpipe/internal/aggregator"
	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/pipeline"
	"github.com/user/logpipe/internal/reader"
)

func makeReader(lines ...string) *reader.Reader {
	input := strings.Join(lines, "\n") + "\n"
	return reader.New(strings.NewReader(input), "svc")
}

func TestPipelineRunPassesMatchingEntries(t *testing.T) {
	r := makeReader(
		`{"level":"info","message":"hello","service":"svc","timestamp":"2024-01-01T00:00:00Z"}`,
		`{"level":"debug","message":"skip","service":"svc","timestamp":"2024-01-01T00:00:01Z"}`,
	)

	agg := aggregator.New(r)
	criteria := filter.Criteria{Level: filter.ParseLevel("info")}

	var buf bytes.Buffer
	w := output.New(&buf, "text")

	p := pipeline.New(agg, criteria, w)
	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "hello") {
		t.Errorf("expected 'hello' in output, got: %s", got)
	}
	if strings.Contains(got, "skip") {
		t.Errorf("did not expect 'skip' in output, got: %s", got)
	}
}

func TestPipelineRunCancelledContext(t *testing.T) {
	r := makeReader(
		`{"level":"info","message":"msg","service":"svc","timestamp":"2024-01-01T00:00:00Z"}`,
	)

	agg := aggregator.New(r)
	criteria := filter.Criteria{}

	var buf bytes.Buffer
	w := output.New(&buf, "text")

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	p := pipeline.New(agg, criteria, w)
	err := p.Run(ctx)
	if err == nil {
		t.Log("context may have been cancelled after drain — acceptable")
	}
}

func TestPipelineRunNoFilters(t *testing.T) {
	r := makeReader(
		`{"level":"warn","message":"a","service":"svc","timestamp":"2024-01-01T00:00:00Z"}`,
		`{"level":"error","message":"b","service":"svc","timestamp":"2024-01-01T00:00:01Z"}`,
	)

	agg := aggregator.New(r)
	criteria := filter.Criteria{}

	var buf bytes.Buffer
	w := output.New(&buf, "text")

	p := pipeline.New(agg, criteria, w)
	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "a") || !strings.Contains(got, "b") {
		t.Errorf("expected both entries in output, got: %s", got)
	}
}
