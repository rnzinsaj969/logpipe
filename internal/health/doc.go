// Package health provides a lightweight health-monitoring layer for logpipe
// log sources.
//
// A Monitor tracks per-source status by recording successes and errors as
// events occur during log ingestion. Callers should invoke RecordSuccess each
// time a log entry is read successfully and RecordError whenever reading or
// parsing fails. After five consecutive errors the source transitions to the
// "down" state; fewer errors result in "degraded".
//
// Snapshots are safe to call from multiple goroutines and return a detached
// copy of the current state so callers can inspect or serialise health data
// without holding any internal locks.
package health
