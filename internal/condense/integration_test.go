package condense_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/logpipe/internal/condense"
	"github.com/your-org/logpipe/internal/reader"
)

func makeEntry(svc, msg, level string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: msg, Level: level}
}

func TestCondenseStreamReducesNoise(t *testing.T) {
	now := time.Unix(2000, 0)
	c, err := condense.New(condense.Options{MinPrefix: 6, MaxAge: 3 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	// inject fixed clock via exported field not available — use Flush timing
	_ = now

	inputs := []reader.LogEntry{
		makeEntry("api", "request handled path=/a", "info"),
		makeEntry("api", "request handled path=/b", "info"),
		makeEntry("api", "request handled path=/c", "info"),
		makeEntry("api", "completely different log line", "warn"),
	}

	var collected []reader.LogEntry
	for _, e := range inputs {
		if out := c.Apply(e); out != nil {
			collected = append(collected, *out)
		}
	}
	for _, e := range c.Flush() {
		collected = append(collected, e)
	}

	// We expect the first three to be condensed (flushed when 4th breaks prefix)
	// and the 4th to remain open until Flush.
	if len(collected) < 2 {
		t.Fatalf("expected at least 2 collected entries, got %d", len(collected))
	}

	// The condensed entry should reference the shared prefix and count.
	found := false
	for _, e := range collected {
		if strings.Contains(e.Message, "3x") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected condensed entry with (3x) count; got: %v", collected)
	}
}

func TestCondenseDoesNotMutateInput(t *testing.T) {
	c, _ := condense.New(condense.DefaultOptions())
	orig := makeEntry("svc", "original message text", "debug")
	copy := orig
	c.Apply(orig)
	if orig.Message != copy.Message || orig.Service != copy.Service {
		t.Fatal("Apply mutated the input entry")
	}
}
