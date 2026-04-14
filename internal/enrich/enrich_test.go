package enrich_test

import (
	"testing"
	"time"

	"github.com/logpipe/internal/enrich"
)

type logEntry struct {
	Timestamp time.Time
	Service   string
	Level     string
	Message   string
	Fields    map[string]string
}

func baseEntry() enrich.Entry {
	return enrich.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Service:   "api",
		Level:     "info",
		Message:   "request received",
		Fields:    map[string]string{"path": "/health"},
	}
}

func TestEnrichAddsStaticField(t *testing.T) {
	e := enrich.New(enrich.WithStaticField("env", "production"))
	out, err := e.Apply(baseEntry())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Fields["env"] != "production" {
		t.Errorf("expected env=production, got %q", out.Fields["env"])
	}
}

func TestEnrichAddsMultipleStaticFields(t *testing.T) {
	e := enrich.New(
		enrich.WithStaticField("env", "staging"),
		enrich.WithStaticField("region", "us-east-1"),
	)
	out, err := e.Apply(baseEntry())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Fields["env"] != "staging" {
		t.Errorf("expected env=staging, got %q", out.Fields["env"])
	}
	if out.Fields["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", out.Fields["region"])
	}
}

func TestEnrichDoesNotMutateOriginal(t *testing.T) {
	e := enrich.New(enrich.WithStaticField("injected", "yes"))
	orig := baseEntry()
	_, err := e.Apply(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := orig.Fields["injected"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestEnrichServiceOverride(t *testing.T) {
	e := enrich.New(enrich.WithServiceOverride("gateway"))
	out, err := e.Apply(baseEntry())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Service != "gateway" {
		t.Errorf("expected service=gateway, got %q", out.Service)
	}
}

func TestEnrichNoOptionsIsNoop(t *testing.T) {
	e := enrich.New()
	in := baseEntry()
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Message != in.Message || out.Service != in.Service || out.Level != in.Level {
		t.Error("noop enricher changed entry")
	}
}

func TestEnrichPreservesExistingFields(t *testing.T) {
	e := enrich.New(enrich.WithStaticField("env", "dev"))
	in := baseEntry()
	in.Fields["path"] = "/api/v1"
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Fields["path"] != "/api/v1" {
		t.Errorf("expected path=/api/v1, got %q", out.Fields["path"])
	}
}
