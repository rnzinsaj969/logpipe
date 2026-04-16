// Package timeslot groups log entries into fixed-duration time buckets.
//
// A Bucketer accepts LogEntry values and assigns each to a Slot whose
// start time is determined by truncating the entry's timestamp to the
// configured slot size (e.g. one minute or one hour).
//
// Typical usage:
//
//	b, err := timeslot.New(timeslot.Options{Size: time.Minute})
//	if err != nil { ... }
//	b.Add(entry)
//	slots := b.Snapshot() // collect and reset
package timeslot
