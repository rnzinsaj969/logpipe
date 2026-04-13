// Package tail provides a file-tailing utility that watches a log file for
// newly appended lines and emits them over a channel.
//
// It is designed to integrate with the logpipe reader pipeline, allowing
// logpipe to consume live log output from services that write to files rather
// than stdout.
//
// Basic usage:
//
//	tailer, err := tail.New("/var/log/myservice/app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer tailer.Close()
//
//	for line := range tailer.Lines(ctx) {
//		fmt.Println(line)
//	}
package tail
