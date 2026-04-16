package alert

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule defines a condition that triggers an alert.
type Rule struct {
	Name     string
	Level    string
	Service  string
	Pattern  string
	compiled *regexp.Regexp
}

// Alert represents a triggered alert event.
type Alert struct {
	Rule  string
	Entry reader.LogEntry
}

// Alerter evaluates log entries against a set of rules and emits alerts.
type Alerter struct {
	mu    sync.Mutex
	rules []Rule
	sink  chan Alert
}

// New creates an Alerter from the given rules. Returns an error if any
// rule contains an invalid pattern.
func New(rules []Rule) (*Alerter, error) {
	compiled := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("alert: rule %q invalid pattern: %w", r.Name, err)
			}
			r.compiled = re
		}
		compiled = append(compiled, r)
	}
	return &Alerter{
		rules: compiled,
		sink:  make(chan Alert, 64),
	}, nil
}

// Evaluate checks entry against all rules and sends matching alerts to the
// internal channel. Non-blocking: alerts are dropped if the buffer is full.
func (a *Alerter) Evaluate(entry reader.LogEntry) {
	a.mu.Lock()
	rules := a.rules
	a.mu.Unlock()

	for _, r := range rules {
		if matches(r, entry) {
			select {
			case a.sink <- Alert{Rule: r.Name, Entry: entry}:
			default:
			}
		}
	}
}

// Alerts returns the channel on which triggered alerts are delivered.
func (a *Alerter) Alerts() <-chan Alert {
	return a.sink
}

func matches(r Rule, e reader.LogEntry) bool {
	if r.Level != "" && r.Level != e.Level {
		return false
	}
	if r.Service != "" && r.Service != e.Service {
		return false
	}
	if r.compiled != nil && !r.compiled.MatchString(e.Message) {
		return false
	}
	return true
}
