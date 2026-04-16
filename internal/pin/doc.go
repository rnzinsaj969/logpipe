// Package pin promotes extra (unstructured) fields on a LogEntry to
// first-class top-level fields — message, level, or service — based on
// a set of user-defined rules.
//
// This is useful when upstream log producers embed the canonical level or
// message inside a nested key rather than the standard fields expected by
// logpipe's pipeline.
//
// Usage:
//
//	p, err := pin.New([]pin.Rule{
//	    {Key: "log_level", Target: "level"},
//	    {Key: "body",      Target: "message"},
//	})
//	if err != nil { ... }
//	out := p.Apply(entry)
package pin
