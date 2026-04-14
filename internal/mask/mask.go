// Package mask provides field-level masking for log entries,
// replacing sensitive field values with a fixed placeholder string.
package mask

import (
	"fmt"

	"github.com/your-org/logpipe/internal/reader"
)

const defaultPlaceholder = "***"

// Options configures the Masker.
type Options struct {
	// Fields is the set of top-level LogEntry field names to mask.
	// Supported targets: "message", "service", and any key in Extra.
	Fields []string
	// Placeholder replaces the original value. Defaults to "***".
	Placeholder string
}

// Masker replaces the values of nominated fields with a placeholder.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
}

// New returns a Masker configured by opts.
// It returns an error if no fields are specified.
func New(opts Options) (*Masker, error) {
	if len(opts.Fields) == 0 {
		return nil, fmt.Errorf("mask: at least one field must be specified")
	}
	ph := opts.Placeholder
	if ph == "" {
		ph = defaultPlaceholder
	}
	fm := make(map[string]struct{}, len(opts.Fields))
	for _, f := range opts.Fields {
		fm[f] = struct{}{}
	}
	return &Masker{fields: fm, placeholder: ph}, nil
}

// Apply returns a copy of entry with nominated fields replaced by the
// placeholder. The original entry is never mutated.
func (m *Masker) Apply(entry reader.LogEntry) reader.LogEntry {
	out := entry

	if _, ok := m.fields["message"]; ok {
		out.Message = m.placeholder
	}
	if _, ok := m.fields["service"]; ok {
		out.Service = m.placeholder
	}

	// Mask extra fields — copy the map to avoid mutating the original.
	for k := range m.fields {
		if k == "message" || k == "service" {
			continue
		}
		if entry.Extra == nil {
			continue
		}
		if _, exists := entry.Extra[k]; exists {
			if out.Extra == nil {
				out.Extra = copyExtra(entry.Extra)
			}
			out.Extra[k] = m.placeholder
		}
	}
	return out
}

// HasField reports whether field is registered for masking.
func (m *Masker) HasField(field string) bool {
	_, ok := m.fields[field]
	return ok
}

func copyExtra(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
