// Package sequence provides a composable, ordered processing chain for
// log entries.
//
// A Sequence is constructed from one or more Processor values. When
// Apply is called the entry flows through each processor in the order
// they were registered. Any processor may transform the entry freely;
// it receives the output of the previous step rather than the original
// input.
//
// If a processor returns an error the chain is aborted immediately and
// the error is wrapped with the failing step index before being
// returned to the caller.
//
// Usage:
//
//	seq, err := sequence.New(normalize, redact, enrich)
//	if err != nil { ... }
//	out, err := seq.Apply(entry)
package sequence
