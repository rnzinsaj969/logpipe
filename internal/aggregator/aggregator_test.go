package aggregator_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/user/logpipe/internal/aggregator"
	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/reader"
)

func makeSource(t *testing.T, name, lines string) aggregator.Source {
	t.Helper()
	r := reader.New(strings.NewReader(lines), name)
	return aggregator.Source{Name: name, Reader: r}
}

func TestAggregatorPassesMatchingEntries(t *testing.T) {
	src := makeSource(t, "svc",
		`{"level":"info","message":"started"}`+"\n"+
			`{"level":"error","message":"boom"}`+"\n")

	var buf bytes.Buffer
	w := output.New(&buf, output.FormatText)
	criteria := filter.Criteria{Level: filter.LevelError}

	agg := aggregator.New([]aggregator.Source{src}, criteria, w)
	if err := agg.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "boom") {
		t.Errorf("expected 'boom' in output, got: %s", got)
	}
	if strings.Contains(got, "started") {
		t.Errorf("'started' should have been filtered out, got: %s", got)
	}
}

func TestAggregatorMultipleSources(t *testing.T) {
	src1 := makeSource(t, "alpha", `{"level":"info","message":"alpha-msg"}`+"\n")
	src2 := makeSource(t, "beta", `{"level":"info","message":"beta-msg"}`+"\n")

	var buf bytes.Buffer
	w := output.New(&buf, output.FormatText)
	agg := aggregator.New([]aggregator.Source{src1, src2}, filter.Criteria{}, w)
	if err := agg.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "alpha-msg") {
		t.Errorf("missing alpha-msg in output: %s", got)
	}
	if !strings.Contains(got, "beta-msg") {
		t.Errorf("missing beta-msg in output: %s", got)
	}
}

func TestAggregatorEmptySource(t *testing.T) {
	src := makeSource(t, "empty", "")
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatText)
	agg := aggregator.New([]aggregator.Source{src}, filter.Criteria{}, w)
	if err := agg.Run(); err != nil && err != io.EOF {
		t.Fatalf("unexpected error: %v", err)
	}
}
