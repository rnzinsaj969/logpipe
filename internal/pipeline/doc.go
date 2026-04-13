// Package pipeline provides the top-level processing loop for logpipe.
//
// A Pipeline combines three components:
//
//   - An [aggregator.Aggregator] that multiplexes log entries from one or more
//     source readers into a single ordered stream.
//
//   - A [filter.Criteria] that decides which entries are forwarded to output.
//
//   - An [output.Writer] that serialises matching entries to the configured
//     destination (stdout, file, etc.) in the requested format.
//
// Typical usage:
//
//	agg := aggregator.New(r1, r2)
//	criteria := filter.Criteria{Level: filter.ParseLevel("warn")}
//	w := output.New(os.Stdout, "json")
//	p := pipeline.New(agg, criteria, w)
//	if err := p.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package pipeline
