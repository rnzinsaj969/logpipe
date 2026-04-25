// Package cast provides lightweight type-coercion helpers for values
// stored in the Extra map of a reader.LogEntry.
//
// Each helper follows the same two-return convention used throughout
// logpipe: (value, ok). A false ok means the key was absent or the
// stored value could not be converted to the requested type.
//
// Supported conversions:
//
//	  String  – returns any scalar as its string representation.
//	  Float64 – converts float64, int, int64, and numeric strings.
//	  Bool    – converts bool and strings accepted by strconv.ParseBool.
//	  Int64   – converts int64, int, float64 (truncated), and numeric strings.
//
// All helpers are safe to call on a nil or empty map; they will simply
// return the zero value and ok=false.
package cast
