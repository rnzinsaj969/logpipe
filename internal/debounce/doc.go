// Package debounce provides a Debouncer that suppresses repeated log entries
// sharing the same (service, message) key within a configurable quiet window.
//
// The first occurrence of a key is always forwarded. Subsequent occurrences
// within the window are suppressed. Once the window expires the next
// occurrence resets the timer and is forwarded again.
//
// Evict may be called periodically to release memory held by stale keys.
package debounce
