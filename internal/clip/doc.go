// Package clip implements numeric field clamping for log entries.
//
// A Clipper is configured with one or more Rules, each specifying a field
// name and a [Min, Max] range. When Apply is called, any numeric value
// stored in LogEntry.Extra under that field is clamped to the range.
// Non-numeric values and absent fields are passed through unchanged.
//
// Validation
//
// New returns an error if any Rule has Min greater than Max, or if the
// same field name appears more than once across the provided rules.
//
// Example usage:
//
//	c, err := clip.New([]clip.Rule{
//		{Field: "latency_ms", Min: 0, Max: 60_000},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	processed := c.Apply(entry)
package clip
