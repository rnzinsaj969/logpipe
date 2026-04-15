package flatten_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/flatten"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Now(),
	}
}

func TestFlattenNoExtra(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	e := base()
	out := f.Apply(e)
	if len(out.Extra) != 0 {
		t.Fatalf("expected empty extra, got %v", out.Extra)
	}
}

func TestFlattenShallowExtra(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	e := base()
	e.Extra = map[string]any{"key": "value", "num": 42}
	out := f.Apply(e)
	if out.Extra["key"] != "value" {
		t.Fatalf("expected key=value, got %v", out.Extra["key"])
	}
	if out.Extra["num"] != 42 {
		t.Fatalf("expected num=42, got %v", out.Extra["num"])
	}
}

func TestFlattenNestedMap(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	e := base()
	e.Extra = map[string]any{
		"http": map[string]any{
			"method": "GET",
			"status": 200,
		},
	}
	out := f.Apply(e)
	if out.Extra["http.method"] != "GET" {
		t.Fatalf("expected http.method=GET, got %v", out.Extra["http.method"])
	}
	if out.Extra["http.status"] != 200 {
		t.Fatalf("expected http.status=200, got %v", out.Extra["http.status"])
	}
}

func TestFlattenCustomSeparator(t *testing.T) {
	f := flatten.New(flatten.Options{Separator: "_"})
	e := base()
	e.Extra = map[string]any{
		"db": map[string]any{"host": "localhost"},
	}
	out := f.Apply(e)
	if out.Extra["db_host"] != "localhost" {
		t.Fatalf("expected db_host=localhost, got %v", out.Extra)
	}
}

func TestFlattenMaxDepth(t *testing.T) {
	f := flatten.New(flatten.Options{Separator: ".", MaxDepth: 1})
	e := base()
	e.Extra = map[string]any{
		"a": map[string]any{
			"b": map[string]any{"c": "deep"},
		},
	}
	out := f.Apply(e)
	// At MaxDepth=1 the nested map under "a" should be flattened one level
	// but the value of "a.b" should remain a map (not recursed further).
	val, ok := out.Extra["a.b"]
	if !ok {
		t.Fatalf("expected key a.b, got %v", out.Extra)
	}
	if _, isMap := val.(map[string]any); !isMap {
		t.Fatalf("expected a.b to be a map at MaxDepth=1, got %T", val)
	}
}

func TestFlattenDoesNotMutateOriginal(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	e := base()
	e.Extra = map[string]any{"x": map[string]any{"y": 1}}
	f.Apply(e)
	if _, nested := e.Extra["x"].(map[string]any); !nested {
		t.Fatal("original extra was mutated")
	}
}

func TestFlatKey(t *testing.T) {
	got := flatten.FlatKey(".", "http", "request", "method")
	want := "http.request.method"
	if got != want {
		t.Fatalf("FlatKey: want %q got %q", want, got)
	}
}
