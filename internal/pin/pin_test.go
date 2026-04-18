package pin_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/pin"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Service: "svc",
		Level:   "info",
		Message: "hello",
		Timestamp: time.Unix(0, 0),
		Extra: map[string]any{
			"msg_override": "world",
			"lvl":          "warn",
			"keep":         "yes",
		},
	}
}

func TestPinMessage(t *testing.T) {
	p, err := pin.New([]pin.Rule{{Key: "msg_override", Target: "message"}})
	if err != nil {
		t.Fatal(err)
	}
	out := p.Apply(base())
	if out.Message != "world" {
		t.Fatalf("expected world, got %s", out.Message)
	}
	if _, ok := out.Extra["msg_override"]; ok {
		t.Fatal("promoted key should be removed from Extra")
	}
}

func TestPinLevel(t *testing.T) {
	p, _ := pin.New([]pin.Rule{{Key: "lvl", Target: "level"}})
	out := p.Apply(base())
	if out.Level != "warn" {
		t.Fatalf("expected warn, got %s", out.Level)
	}
}

func TestPinService(t *testing.T) {
	entry := base()
	entry.Extra["svc_override"] = "newsvc"
	p, _ := pin.New([]pin.Rule{{Key: "svc_override", Target: "service"}})
	out := p.Apply(entry)
	if out.Service != "newsvc" {
		t.Fatalf("expected newsvc, got %s", out.Service)
	}
	if _, ok := out.Extra["svc_override"]; ok {
		t.Fatal("promoted key should be removed from Extra")
	}
}

func TestPinDoesNotMutateOriginal(t *testing.T) {
	entry := base()
	p, _ := pin.New([]pin.Rule{{Key: "msg_override", Target: "message"}})
	p.Apply(entry)
	if entry.Message != "hello" {
		t.Fatal("original entry was mutated")
	}
	if _, ok := entry.Extra["msg_override"]; !ok {
		t.Fatal("original extra was mutated")
	}
}

func TestPinMissingKeyIsNoOp(t *testing.T) {
	p, _ := pin.New([]pin.Rule{{Key: "absent", Target: "service"}})
	out := p.Apply(base())
	if out.Service != "svc" {
		t.Fatal("service should be unchanged")
	}
}

func TestPinNonStringValueSkipped(t *testing.T) {
	entry := base()
	entry.Extra["num"] = 42
	p, _ := pin.New([]pin.Rule{{Key: "num", Target: "service"}})
	out := p.Apply(entry)
	if out.Service != "svc" {
		t.Fatal("non-string extra should not pin")
	}
}

func TestNewEmptyKeyReturnsError(t *testing.T) {
	_, err := pin.New([]pin.Rule{{Key: "", Target: "level"}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewUnsupportedTargetReturnsError(t *testing.T) {
	_, err := pin.New([]pin.Rule{{Key: "k", Target: "timestamp"}})
	if err == nil {
		t.Fatal("expected error for unsupported target")
	}
}
