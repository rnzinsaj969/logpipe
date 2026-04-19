// Package gate provides a runtime-togglable gate that either passes log
// entries through unchanged or drops them.
//
// A Gate is safe for concurrent use. Call Open and Close from any goroutine
// to change the state; Apply is the per-entry hot path used inside a
// processing pipeline.
//
// Example:
//
//	g, _ := gate.New(true)  // starts open
//	entry, err := g.Apply(e)
//	if err != nil {
//		// entry was dropped
//	}
//	g.Close() // subsequent calls to Apply will return an error
package gate
