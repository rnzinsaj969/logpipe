package route

import (
	"errors"
	"regexp"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule defines a single routing rule: entries matching the predicate are
// forwarded to the named destination.
type Rule struct {
	// Destination is an opaque label identifying the output sink.
	Destination string
	// Level, if non-empty, matches only entries with this exact level.
	Level string
	// ServicePattern, if non-empty, is a regular expression matched against
	// the entry's Service field.
	ServicePattern string

	compiledService *regexp.Regexp
}

// Router dispatches log entries to one or more destinations based on an
// ordered list of rules. The first matching rule wins.
type Router struct {
	rules []Rule
}

// New compiles all ServicePattern fields and returns a ready-to-use Router.
// An error is returned if any pattern fails to compile.
func New(rules []Rule) (*Router, error) {
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		if r.Destination == "" {
			return nil, errors.New("route: rule missing destination")
		}
		if r.ServicePattern != "" {
			re, err := regexp.Compile(r.ServicePattern)
			if err != nil {
				return nil, err
			}
			r.compiledService = re
		}
		compiled[i] = r
	}
	return &Router{rules: compiled}, nil
}

// Match returns the destination label for the first rule that matches entry.
// If no rule matches, an empty string is returned.
func (rt *Router) Match(entry reader.LogEntry) string {
	for _, r := range rt.rules {
		if r.Level != "" && entry.Level != r.Level {
			continue
		}
		if r.compiledService != nil && !r.compiledService.MatchString(entry.Service) {
			continue
		}
		return r.Destination
	}
	return ""
}

// Rules returns a copy of the router's rule list.
func (rt *Router) Rules() []Rule {
	out := make([]Rule, len(rt.rules))
	copy(out, rt.rules)
	return out
}
