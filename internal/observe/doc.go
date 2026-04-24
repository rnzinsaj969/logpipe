// Package observe provides a non-destructive tap stage for a logpipe
// processing pipeline.
//
// An Observer sits inline in a pipeline and invokes a user-supplied Handler
// for every log entry it sees. The entry is forwarded downstream unchanged,
// making it easy to collect metrics, write debug output, or trigger side
// effects without altering the stream.
//
// Basic usage:
//
//	obs, err := observe.New(func(e reader.LogEntry) {
//		fmt.Println(e.Message)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, _ := obs.Apply(entry)
package observe
