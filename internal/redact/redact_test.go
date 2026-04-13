package redact

import (
	"testing"
)

func TestApplyRedactsBearer(t *testing.T) {
	r, err := NewFromPatterns([]string{`(?i)bearer\s+[A-Za-z0-9\-._~+/]+=*`}, DefaultReplacement)
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	input := "Authorization: Bearer abc123token"
	got := r.Apply(input)
	if got == input {
		t.Errorf("expected redaction, got original string")
	}
	if got != "Authorization: "+DefaultReplacement {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestApplyNoMatchReturnsOriginal(t *testing.T) {
	r, err := NewFromPatterns([]string{`password=[^\s&]+`}, DefaultReplacement)
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	input := "level=info msg=hello"
	if got := r.Apply(input); got != input {
		t.Errorf("expected no change, got %q", got)
	}
}

func TestApplyMapRedactsValues(t *testing.T) {
	r, err := NewFromPatterns([]string{`password=[^\s&]+`}, DefaultReplacement)
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	fields := map[string]string{
		"msg":   "login password=secret123",
		"level": "info",
	}
	out := r.ApplyMap(fields)
	if out["level"] != "info" {
		t.Errorf("non-sensitive field altered: %q", out["level"])
	}
	if out["msg"] == fields["msg"] {
		t.Errorf("expected msg to be redacted")
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := NewFromPatterns([]string{`[invalid`}, DefaultReplacement)
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestHasRules(t *testing.T) {
	empty := New(nil)
	if empty.HasRules() {
		t.Error("expected HasRules to be false for empty redactor")
	}
	r, _ := NewFromPatterns([]string{`foo`}, "[X]")
	if !r.HasRules() {
		t.Error("expected HasRules to be true")
	}
}

func TestMask(t *testing.T) {
	cases := []struct {
		input string
		n     int
		want  string
	}{
		{"abcdefgh", 2, "ab****gh"},
		{"ab", 2, "**"},
		{"abcd", 2, "****"},
		{"hello", 1, "h***o"},
	}
	for _, tc := range cases {
		got := Mask(tc.input, tc.n)
		if got != tc.want {
			t.Errorf("Mask(%q, %d) = %q, want %q", tc.input, tc.n, got, tc.want)
		}
	}
}

func TestNewDefault(t *testing.T) {
	r, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault error: %v", err)
	}
	if !r.HasRules() {
		t.Error("expected default redactor to have rules")
	}
}

func TestRedactStringEmail(t *testing.T) {
	input := "contact user@example.com for help"
	out := RedactString(input)
	if out == input {
		t.Errorf("expected email to be redacted in %q", input)
	}
}
