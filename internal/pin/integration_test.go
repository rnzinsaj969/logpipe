package pin_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/pin"
	"github.com/logpipe/logpipe/internal/reader"
)

func TestPinMultipleRulesAppliedInOrder(t *testing.T) {
	rules := []pin.Rule{
		{Key: "msg_override", Target: "message"},
		{Key: "lvl", Target: "level"},
		{Key: "src", Target: "service"},
	}
	p, err := pin.New(rules)
	if err != nil {
		t.Fatal(err)
	}
	entry := reader.LogEntry{
		Service:   "old",
		Level:     "info",
		Message:   "old msg",
		Timestamp: time.Now(),
		Extra: map[string]any{
			"msg_override": "new msg",
			"lvl":          "error",
			"src":          "auth",
		},
	}
	out := p.Apply(entry)
	if out.Message != "new msg" || out.Level != "error" || out.Service != "auth" {
		t.Fatalf("unexpected result: %+v", out)
	}
	if len(out.Extra) != 0 {
		t.Fatalf("expected empty extra, got %v", out.Extra)
	}
}

func TestPinChainedWithExtraFieldsPreserved(t *testing.T) {
	p, _ := pin.New([]pin.Rule{{Key: "lvl", Target: "level"}})
	entry := reader.LogEntry{
		Level:     "info",
		Message:   "msg",
		Timestamp: time.Now(),
		Extra: map[string]any{
			"lvl":  "debug",
			"keep": "retained",
		},
	}
	out := p.Apply(entry)
	if out.Level != "debug" {
		t.Fatalf("expected debug, got %s", out.Level)
	}
	if out.Extra["keep"] != "retained" {
		t.Fatal("unrelated extra field should be preserved")
	}
}
