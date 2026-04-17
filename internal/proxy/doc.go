// Package proxy provides a fan-out forwarder that delivers log entries to one
// or more named Sink implementations. Sinks can be registered and removed at
// runtime. Forward collects and returns all sink errors so the caller can
// decide how to handle partial failures.
package proxy
