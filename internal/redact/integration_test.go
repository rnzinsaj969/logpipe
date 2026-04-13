package redact_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/redact"
)

func TestRedactorChainMultiplePatterns(t *testing.T) {
	patterns := []string{
		`(?i)bearer\s+[A-Za-z0-9\-._~+/]+=*`,
		`password=[^\s&]+`,
	}
	r, err := redact.NewFromPatterns(patterns, redact.DefaultReplacement)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	input := "Authorization: Bearer tok3n password=s3cr3t"
	out := r.Apply(input)

	if strings.Contains(out, "tok3n") {
		t.Errorf("token not redacted: %q", out)
	}
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("password not redacted: %q", out)
	}
}

func TestApplyMapDoesNotMutateOriginal(t *testing.T) {
	r, _ := redact.NewFromPatterns([]string{`password=[^\s&]+`}, redact.DefaultReplacement)
	original := map[string]string{
		"msg": "login password=hunter2",
	}
	out := r.ApplyMap(original)
	if original["msg"] != "login password=hunter2" {
		t.Error("ApplyMap mutated the original map")
	}
	if out["msg"] == original["msg"] {
		t.Error("output map was not redacted")
	}
}

func TestRedactStringCreditCard(t *testing.T) {
	input := "card 4111 1111 1111 1111 was charged"
	out := redact.RedactString(input)
	if strings.Contains(out, "4111") {
		t.Errorf("credit card number not redacted: %q", out)
	}
}
