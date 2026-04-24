// Package extract implements a log entry processor that promotes values
// from the Extra map into the top-level fields of a LogEntry.
//
// This is useful when upstream services embed the canonical message, level,
// or service name inside a structured field rather than the standard keys.
//
// Example usage:
//
//	e, err := extract.New([]extract.Rule{
//		{From: "log_message", To: extract.TargetMessage, Overwrite: true},
//		{From: "severity",   To: extract.TargetLevel,   Overwrite: false},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	output := e.Apply(entry)
package extract
