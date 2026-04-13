// Package redact provides pattern-based redaction of sensitive data in log
// messages before they are written to any output sink.
//
// A Redactor holds one or more compiled regular-expression rules. Each rule
// specifies a pattern to match and a replacement string (typically a fixed
// token such as "[REDACTED]"). Rules are applied in order; the output of one
// rule becomes the input of the next.
//
// Usage:
//
//	r, err := redact.NewFromPatterns([]string{`password=[^\s&]+`}, "[REDACTED]")
//	if err != nil { ... }
//	clean := r.Apply(rawMessage)
//
// A set of DefaultPatterns covering bearer tokens, passwords, credit-card
// numbers, and email addresses is available via NewDefault.
package redact
