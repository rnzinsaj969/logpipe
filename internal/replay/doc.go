// Package replay provides a Replayer that re-emits a captured sequence of
// log entries onto a channel, honouring the original inter-entry timing
// scaled by a configurable speed multiplier.
//
// Typical use-cases include:
//   - Replaying historical log data through the processing pipeline during
//     testing or debugging.
//   - Stress-testing downstream consumers at accelerated speed.
//
// Example:
//
//	r, err := replay.New(entries, replay.Options{Speed: 2.0})
//	if err != nil { ... }
//	out := make(chan reader.LogEntry, 64)
//	go r.Run(ctx, out)
package replay
