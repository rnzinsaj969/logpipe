// Package field provides a Processor that applies copy, rename, and delete
// operations to the Extra fields of a log entry.
//
// Operations are applied in order. Mutations do not affect the original entry.
//
// Example usage:
//
//	op, err := field.New([]field.Op{
//		{Action: "rename", From: "request_id", To: "req_id"},
//		{Action: "delete", From: "internal_token"},
//	})
//	if err != nil { ... }
//	result := op.Apply(entry)
package field
