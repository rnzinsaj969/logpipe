package annotate

import (
	"errors"
	"regexp"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule describes a single annotation rule.
type Rule struct {
	// Pattern is a regular expression matched against the log message.
	Pattern string
	// Key is the extra field key to set when the rule matches.
	Key string
	// Value is the extra field value to set.
	Value string

	re *regexp.Regexp
}

// Annotator applies annotation rules to log entries.
type Annotator struct {
	rules []Rule
}

// New returns an Annotator for the given rules.
// Returns an error if any pattern fails to compile or a rule is incomplete.
func New(rules []Rule) (*Annotator, error) {
	if len(rules) == 0 {
		return nil, errors.New("annotate: at least one rule is required")
	}
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		if r.Key == "" {
			return nil, errors.New("annotate: rule key must not be empty")
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, errors.New("annotate: invalid pattern: " + err.Error())
		}
		compiled[i] = r
		compiled[i].re = re
	}
	return &Annotator{rules: compiled}, nil
}

// Apply evaluates all rules against e and returns a new entry with any
// matching annotations added to Extra. The original entry is not mutated.
func (a *Annotator) Apply(e reader.LogEntry) reader.LogEntry {
	out := e
	out.Extra = copyExtra(e.Extra)
	for _, r := range a.rules {
		if r.re.MatchString(e.Message) {
			out.Extra[r.Key] = r.Value
		}
	}
	return out
}

func copyExtra(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
