# logpipe

A lightweight CLI tool for structured log aggregation and real-time filtering across multiple services.

---

## Installation

```bash
go install github.com/yourusername/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logpipe.git && cd logpipe && go build -o logpipe .
```

---

## Usage

Pipe logs from multiple services and filter in real time:

```bash
# Aggregate logs from multiple sources and filter by level
logpipe --sources api.log,worker.log --filter level=error

# Stream logs from a running service with JSON output
logpipe --follow --source /var/log/myapp.log --format json

# Filter by service name and minimum log level
logpipe --sources api.log,db.log --service api --level warn
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--sources` | Comma-separated list of log files or stdin | `stdin` |
| `--filter` | Key=value filter expression | none |
| `--level` | Minimum log level (debug, info, warn, error) | `info` |
| `--follow` | Stream logs in real time | `false` |
| `--format` | Output format: `json` or `text` | `text` |

---

## Features

- Real-time log streaming with `--follow`
- Structured JSON log parsing and filtering
- Aggregation across multiple log sources simultaneously
- Lightweight with no external dependencies

---

## License

MIT © 2024 yourusername