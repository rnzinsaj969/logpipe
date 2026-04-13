package redact

import (
	"fmt"
	"testing"
)

func BenchmarkApplyDefault(b *testing.B) {
	r, err := NewDefault()
	if err != nil {
		b.Fatalf("setup: %v", err)
	}
	input := "user user@example.com logged in with Bearer eyJhbGciOiJIUzI1NiJ9.payload.sig"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Apply(input)
	}
}

func BenchmarkApplyMapSmall(b *testing.B) {
	r, err := NewDefault()
	if err != nil {
		b.Fatalf("setup: %v", err)
	}
	fields := map[string]string{
		"msg":     "contact admin@corp.io for access",
		"level":   "warn",
		"service": "auth",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.ApplyMap(fields)
	}
}

func BenchmarkMask(b *testing.B) {
	s := fmt.Sprintf("%032d", 12345678901234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Mask(s, 4)
	}
}
