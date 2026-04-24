// Package signal implements a lightweight in-process publish/subscribe bus
// for broadcasting named events across logpipe components.
//
// A Bus allows any number of handlers to subscribe to a named signal.
// When Fire is called, all registered handlers receive the signal name
// and an optional payload map. Subscriptions are safe for concurrent use.
//
// Example:
//
//	bus := signal.New()
//	_ = bus.Subscribe("log.error", func(name string, payload map[string]any) {
//		fmt.Println("signal fired:", name, payload)
//	})
//	bus.Fire("log.error", map[string]any{"source": "api"})
package signal
