// Package cast provides type-coercion helpers for log entry fields.
// It converts extra-field values to common scalar types, returning a
// typed result and a boolean indicating whether the conversion succeeded.
package cast

import (
	"fmt"
	"strconv"

	"github.com/logpipe/logpipe/internal/reader"
)

// String returns the named extra field as a string.
// Numeric and boolean values are formatted with fmt.Sprintf.
// Missing keys return ("", false).
func String(e reader.LogEntry, key string) (string, bool) {
	v, ok := e.Extra[key]
	if !ok {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	case bool:
		return strconv.FormatBool(t), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

// Float64 returns the named extra field as a float64.
// Strings are parsed with strconv.ParseFloat.
// Missing keys or unconvertible values return (0, false).
func Float64(e reader.LogEntry, key string) (float64, bool) {
	v, ok := e.Extra[key]
	if !ok {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}

// Bool returns the named extra field as a bool.
// String values are parsed with strconv.ParseBool.
// Missing keys or unconvertible values return (false, false).
func Bool(e reader.LogEntry, key string) (bool, bool) {
	v, ok := e.Extra[key]
	if !ok {
		return false, false
	}
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		b, err := strconv.ParseBool(t)
		if err != nil {
			return false, false
		}
		return b, true
	}
	return false, false
}
