// Package window provides a thread-safe sliding-window event counter.
//
// A Counter records the timestamps of events as they are added and
// automatically evicts entries that fall outside the configured window
// duration. This makes it straightforward to answer questions such as
// "how many log lines arrived in the last 10 seconds?" without
// accumulating unbounded memory.
//
// Usage:
//
//	c := window.New(10 * time.Second)
//	c.Add()
//	fmt.Println(c.Count()) // 1
package window
