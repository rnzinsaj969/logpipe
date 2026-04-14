package parse_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/parse"
)

func TestStringPresentStringValue(t *testing.T) {
	f := parse.Fields{"env": "production"}
	v, ok := parse.String(f, "env")
	if !ok || v != "production" {
		t.Fatalf("expected production/true, got %q/%v", v, ok)
	}
}

func TestStringMissingKey(t *testing.T) {
	f := parse.Fields{}
	_, ok := parse.String(f, "missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestStringNonStringValue(t *testing.T) {
	f := parse.Fields{"code": 42}
	v, ok := parse.String(f, "code")
	if !ok || v != "42" {
		t.Fatalf("expected \"42\"/true, got %q/%v", v, ok)
	}
}

func TestIntFromInt(t *testing.T) {
	f := parse.Fields{"count": 7}
	n, err := parse.Int(f, "count")
	if err != nil || n != 7 {
		t.Fatalf("expected 7/nil, got %d/%v", n, err)
	}
}

func TestIntFromFloat64(t *testing.T) {
	f := parse.Fields{"score": float64(3.9)}
	n, err := parse.Int(f, "score")
	if err != nil || n != 3 {
		t.Fatalf("expected 3/nil, got %d/%v", n, err)
	}
}

func TestIntFromString(t *testing.T) {
	f := parse.Fields{"port": "8080"}
	n, err := parse.Int(f, "port")
	if err != nil || n != 8080 {
		t.Fatalf("expected 8080/nil, got %d/%v", n, err)
	}
}

func TestIntMissingKey(t *testing.T) {
	_, err := parse.Int(parse.Fields{}, "x")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestIntInvalidString(t *testing.T) {
	f := parse.Fields{"val": "notanumber"}
	_, err := parse.Int(f, "val")
	if err == nil {
		t.Fatal("expected error for non-numeric string")
	}
}

func TestBoolFromBool(t *testing.T) {
	f := parse.Fields{"enabled": true}
	b, err := parse.Bool(f, "enabled")
	if err != nil || !b {
		t.Fatalf("expected true/nil, got %v/%v", b, err)
	}
}

func TestBoolFromString(t *testing.T) {
	f := parse.Fields{"debug": "false"}
	b, err := parse.Bool(f, "debug")
	if err != nil || b {
		t.Fatalf("expected false/nil, got %v/%v", b, err)
	}
}

func TestBoolMissingKey(t *testing.T) {
	_, err := parse.Bool(parse.Fields{}, "flag")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
