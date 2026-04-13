// Package aggregator provides fan-in log aggregation across multiple
// service readers.
//
// Usage:
//
//	sources := []aggregator.Source{
//		{Name: "api",    Reader: reader.New(apiStream,    "api")},
//		{Name: "worker", Reader: reader.New(workerStream, "worker")},
//	}
//
//	criteria := filter.Criteria{Level: filter.LevelWarn}
//	out     := output.New(os.Stdout, output.FormatText)
//
//	agg := aggregator.New(sources, criteria, out)
//	if err := agg.Run(); err != nil {
//		log.Fatal(err)
//	}
//
// Each source is consumed in its own goroutine. Entries that satisfy the
// filter criteria are forwarded to the shared output writer in the order
// they arrive.
package aggregator
