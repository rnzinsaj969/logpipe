package cast_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/cast"
	"github.com/logpipe/logpipe/internal/reader"
)

func base(extra map[string]any) reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "msg",
		Timestamp: time.Now(),
		Extra:     extra,
	}
}

// String tests

func TestStringFromString(t *testing.T) {
	e := base(map[string]any{"k": "hello"})
	v, ok := cast.String(e, "k")
	if !ok || v != "hello" {
		t.Fatalf("want hello/true, got %q/%v", v, ok)
	}
}

func TestStringFromFloat(t *testing.T) {
	e := base(map[string]any{"k": float64(3.14)})
	v, ok := cast.String(e, "k")
	if !ok || v == "" {
		t.Fatalf("want non-empty/true, got %q/%v", v, ok)
	}
}

func TestStringFromBool(t *testing.T) {
	e := base(map[string]any{"k": true})
	v, ok := cast.String(e, "k")
	if !ok || v != "true" {
		t.Fatalf("want true/true, got %q/%v", v, ok)
	}
}

func TestStringMissingKey(t *testing.T) {
	e := base(map[string]any{})
	v, ok := cast.String(e, "missing")
	if ok || v != "" {
		t.Fatalf("want empty/false, got %q/%v", v, ok)
	}
}

// Float64 tests

func TestFloat64FromFloat(t *testing.T) {
	e := base(map[string]any{"n": float64(2.5)})
	v, ok := cast.Float64(e, "n")
	if !ok || v != 2.5 {
		t.Fatalf("want 2.5/true, got %v/%v", v, ok)
	}
}

func TestFloat64FromString(t *testing.T) {
	e := base(map[string]any{"n": "1.23"})
	v, ok := cast.Float64(e, "n")
	if !ok || v != 1.23 {
		t.Fatalf("want 1.23/true, got %v/%v", v, ok)
	}
}

func TestFloat64FromInt(t *testing.T) {
	e := base(map[string]any{"n": 7})
	v, ok := cast.Float64(e, "n")
	if !ok || v != 7.0 {
		t.Fatalf("want 7/true, got %v/%v", v, ok)
	}
}

func TestFloat64BadString(t *testing.T) {
	e := base(map[string]any{"n": "abc"})
	_, ok := cast.Float64(e, "n")
	if ok {
		t.Fatal("want false for non-numeric string")
	}
}

func TestFloat64MissingKey(t *testing.T) {
	e := base(map[string]any{})
	_, ok := cast.Float64(e, "missing")
	if ok {
		t.Fatal("want false for missing key")
	}
}

// Bool tests

func TestBoolFromBool(t *testing.T) {
	e := base(map[string]any{"b": false})
	v, ok := cast.Bool(e, "b")
	if !ok || v != false {
		t.Fatalf("want false/true, got %v/%v", v, ok)
	}
}

func TestBoolFromString(t *testing.T) {
	e := base(map[string]any{"b": "true"})
	v, ok := cast.Bool(e, "b")
	if !ok || !v {
		t.Fatalf("want true/true, got %v/%v", v, ok)
	}
}

func TestBoolBadString(t *testing.T) {
	e := base(map[string]any{"b": "yes"})
	_, ok := cast.Bool(e, "b")
	if ok {
		t.Fatal("want false for unparseable bool string")
	}
}

func TestBoolMissingKey(t *testing.T) {
	e := base(map[string]any{})
	_, ok := cast.Bool(e, "missing")
	if ok {
		t.Fatal("want false for missing key")
	}
}
