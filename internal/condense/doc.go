// Package condense provides a Condenser that merges consecutive log entries
// from the same service when their messages share a common prefix.
//
// This reduces noise from services that emit many near-identical lines (e.g.
// repeated retries or partial-line flushes) by collapsing them into a single
// representative entry annotated with a repeat count.
//
// Basic usage:
//
//	c, err := condense.New(condense.DefaultOptions())
//	if err != nil { ... }
//
//	for _, e := range entries {
//		if out := c.Apply(e); out != nil {
//			process(*out)
//		}
//	}
//	// flush remaining open groups
//	for _, e := range c.Flush() {
//		process(e)
//	}
package condense
