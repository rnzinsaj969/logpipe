// Package group provides a Grouper that accumulates log entries into
// named buckets based on a configurable field (level, service, or any
// Extra key). Snapshots are atomic and reset internal state, making
// the Grouper suitable for periodic aggregation windows.
package group
