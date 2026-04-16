// Package drop provides a rule-based log entry dropper for logpipe.
//
// A Dropper holds a set of Rules; each Rule may specify a Level, a Service
// name, and/or a regular-expression Pattern applied to the log message.
// An entry is dropped when it satisfies ALL non-empty fields of at least
// one Rule.
//
// Example usage:
//
//	d, err := drop.New([]drop.Rule{
//		{Level: "debug"},
//		{Service: "healthd", Pattern: `^ping`},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, e := range entries {
//		if !d.ShouldDrop(e) {
//			process(e)
//		}
//	}
package drop
