// Package parse provides utilities for extracting and coercing typed
// values from the extra fields map attached to a log entry.
package parse

import (
	"fmt"
	"strconv"
)

// Fields is an alias for the extra-fields map used throughout logpipe.
type Fields = map[string]any

// String returns the value at key as a string.
// If the key is absent the empty string and false are returned.
// If the value is not a string it is formatted with fmt.Sprintf.
func String(fields Fields, key string) (string, bool) {
	v, ok := fields[key]
	if !ok {
		return "", false
	}
	if s, ok := v.(string); ok {
		return s, true
	}
	return fmt.Sprintf("%v", v), true
}

// Int returns the value at key coerced to int64.
// Supported source types: int, int64, float64, string.
// Returns 0 and an error when the value cannot be coerced.
func Int(fields Fields, key string) (int64, error) {
	v, ok := fields[key]
	if !ok {
		return 0, fmt.Errorf("parse: key %q not found", key)
	}
	switch t := v.(type) {
	case int:
		return int64(t), nil
	case int64:
		return t, nil
	case float64:
		return int64(t), nil
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse: cannot convert %q to int: %w", t, err)
		}
		return n, nil
	}
	return 0, fmt.Errorf("parse: unsupported type %T for key %q", v, key)
}

// Bool returns the value at key coerced to bool.
// Supported source types: bool, string ("true"/"false", case-insensitive).
func Bool(fields Fields, key string) (bool, error) {
	v, ok := fields[key]
	if !ok {
		return false, fmt.Errorf("parse: key %q not found", key)
	}
	switch t := v.(type) {
	case bool:
		return t, nil
	case string:
		b, err := strconv.ParseBool(t)
		if err != nil {
			return false, fmt.Errorf("parse: cannot convert %q to bool: %w", t, err)
		}
		return b, nil
	}
	return false, fmt.Errorf("parse: unsupported type %T for key %q", v, key)
}
