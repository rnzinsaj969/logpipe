// Package pivot provides a log entry aggregator that groups entries by a
// specified field and maintains per-group counts.
//
// A Pivot accumulates entries into named buckets determined by a field value
// extracted from each entry. Callers can take a snapshot of the current
// group counts at any time; the snapshot is an isolated copy.
//
// Example usage:
//
//	p, err := pivot.New("level")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	p.Add(entry)
//	counts := p.Snapshot()
package pivot
