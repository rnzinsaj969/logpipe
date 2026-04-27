// Package surge provides a per-service burst detector for log streams.
//
// A Detector maintains a sliding window of event timestamps for each
// service name. On every call to Record the detector compares the
// recent arrival rate (last quarter of the configured window) against
// the long-term baseline rate for that service. When the recent rate
// exceeds the baseline by the configured multiple the call returns true,
// indicating a surge that the caller may wish to alert on or throttle.
//
// Example:
//
//	d, err := surge.New(time.Minute, 3.0)
//	if err != nil { ... }
//	if d.Record(entry) {
//		log.Printf("surge detected for service %s", entry.Service)
//	}
package surge
