# logslice

Stream and filter structured JSON logs from remote servers with regex and time-range support.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build ./...
```

---

## Usage

Stream logs from a remote server and filter by pattern or time range:

```bash
# Filter logs matching a regex pattern
logslice --host logs.example.com --filter "error|panic"

# Filter logs within a time range
logslice --host logs.example.com --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z"

# Combine regex and time range
logslice --host logs.example.com --filter "timeout" --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z"

# Pretty-print matched JSON entries
logslice --host logs.example.com --filter "user_id=42" --pretty
```

### Flags

| Flag | Description |
|------|-------------|
| `--host` | Remote server address |
| `--filter` | Regex pattern to match against log entries |
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--pretty` | Pretty-print JSON output |
| `--field` | JSON field to apply the regex filter against (default: full line) |

---

## Example Output

```json
{"timestamp":"2024-01-15T08:32:11Z","level":"error","message":"connection timeout","service":"api"}
{"timestamp":"2024-01-15T08:45:02Z","level":"error","message":"timeout waiting for db","service":"worker"}
```

---

## License

MIT © 2024 yourname