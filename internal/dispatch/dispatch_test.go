package dispatch_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/dispatch"
)

type mockEntry struct {
	Service string
	Level   string
	Message string
}

func makeEntry(svc, lvl, msg string) dispatch.Entry {
	return dispatch.Entry{
		Service:   svc,
		Level:     lvl,
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func TestDispatchRoutesByKey(t *testing.T) {
	d, err := dispatch.New(func(e dispatch.Entry) string { return e.Level })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var gotError, gotInfo string
	d.Register("error", func(e dispatch.Entry) error { gotError = e.Message; return nil })
	d.Register("info", func(e dispatch.Entry) error { gotInfo = e.Message; return nil })

	if err := d.Dispatch(makeEntry("svc", "error", "boom")); err != nil {
		t.Fatalf("dispatch error: %v", err)
	}
	if err := d.Dispatch(makeEntry("svc", "info", "hello")); err != nil {
		t.Fatalf("dispatch info: %v", err)
	}

	if gotError != "boom" {
		t.Errorf("error sink: got %q, want %q", gotError, "boom")
	}
	if gotInfo != "hello" {
		t.Errorf("info sink: got %q, want %q", gotInfo, "hello")
	}
}

func TestDispatchUnknownKeyIsNoop(t *testing.T) {
	d, err := dispatch.New(func(e dispatch.Entry) string { return e.Level })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Dispatch(makeEntry("svc", "warn", "msg")); err != nil {
		t.Errorf("expected nil for unknown key, got %v", err)
	}
}

func TestDispatchNilKeyFuncReturnsError(t *testing.T) {
	_, err := dispatch.New(nil)
	if err == nil {
		t.Fatal("expected error for nil key func")
	}
}

func TestDispatchSinkError(t *testing.T) {
	d, err := dispatch.New(func(e dispatch.Entry) string { return e.Level })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentinel := errors.New("sink failure")
	d.Register("error", func(e dispatch.Entry) error { return sentinel })

	if err := d.Dispatch(makeEntry("svc", "error", "msg")); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestDispatchRegisterOverwrites(t *testing.T) {
	d, err := dispatch.New(func(e dispatch.Entry) string { return e.Level })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var calls int
	d.Register("info", func(e dispatch.Entry) error { calls++; return nil })
	d.Register("info", func(e dispatch.Entry) error { calls += 10; return nil })

	if err := d.Dispatch(makeEntry("svc", "info", "msg")); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if calls != 10 {
		t.Errorf("expected second handler to replace first, calls=%d", calls)
	}
}
