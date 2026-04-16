// Package alert provides rule-based alerting for log entries.
//
// An Alerter holds a set of Rules, each describing conditions on the
// Level, Service, and/or Message of a log entry. When Evaluate is called
// with a LogEntry, any rule whose conditions all match causes an Alert to
// be emitted on the channel returned by Alerts.
//
// Rules are evaluated independently, so a single entry may trigger
// multiple alerts. Alert delivery is non-blocking: if the internal buffer
// is full, excess alerts are silently dropped.
//
// Example:
//
//	rules := []alert.Rule{
//		{Name: "high-error-rate", Level: "error", Service: "api"},
//		{Name: "panic-detected",  Pattern: `panic`},
//	}
//	a, err := alert.New(rules)
//	if err != nil { ... }
//	for entry := range entries {
//		a.Evaluate(entry)
//	}
package alert
