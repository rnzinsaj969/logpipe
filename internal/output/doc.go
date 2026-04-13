// Package output provides formatting and writing capabilities for structured
// log entries produced by logpipe.
//
// It supports two output formats:
//
//   - text: a human-readable single-line format suitable for terminal output,
//     e.g. "2024-01-15T10:30:00Z [INFO] (api) request received"
//
//   - json: a machine-readable JSON format where each entry is emitted as a
//     single newline-delimited JSON object, suitable for piping into other tools.
//
// Example usage:
//
//	w := output.New(os.Stdout, output.FormatText)
//	w.Write(output.Entry{
//		Timestamp: time.Now(),
//		Service:   "api",
//		Level:     "INFO",
//		Message:   "server started",
//	})
package output
