// Package coalesce provides a field-coalescing processor for log entries.
//
// A Coalescer is configured with one or more Rules. Each Rule specifies an
// ordered list of source Extra field keys and a single target key. When Apply
// is called the first source field that contains a non-empty string value is
// written to the target field, leaving all other fields untouched.
//
// This is useful when different upstream services emit the same semantic value
// under different key names (e.g. "host", "hostname", "node") and the
// pipeline needs a single canonical field for downstream filtering or output.
//
// Example:
//
//	rules := []coalesce.Rule{
//		{Sources: []string{"host", "hostname", "node"}, Target: "canonical_host"},
//	}
//	c, err := coalesce.New(rules)
//	if err != nil { ... }
//	processed := c.Apply(entry)
package coalesce
