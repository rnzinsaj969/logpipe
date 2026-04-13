// Package buffer provides a thread-safe fixed-capacity ring buffer
// for log entries. It is used to retain a sliding window of recent
// entries in memory, enabling features such as burst replay and
// back-pressure absorption between pipeline stages.
//
// When the buffer reaches capacity, the oldest entry is silently
// discarded to make room for the newest one, ensuring that the
// most recent log activity is always preserved.
package buffer
