// Package tally provides a concurrency-safe frequency counter for log entry fields.
//
// A Tally instance tracks how many times each unique value appears for a
// configured field (e.g. "level", "service", or any Extra key) across the
// entries passed to Add. The current counts can be retrieved at any time via
// Snapshot, which returns an isolated copy, and Reset clears all state.
//
// Example usage:
//
//	t, _ := tally.New("level")
//	for _, e := range entries {
//		t.Add(e)
//	}
//	fmt.Println(t.Snapshot()) // map[error:3 info:17 warn:2]
package tally
