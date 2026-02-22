# CLAUDE.md

This file provides guidance for AI assistants working on this codebase.

## Project Overview

`es` is a Go CLI tool for interacting with Elasticsearch REST APIs from the terminal. It provides commands for indexing documents, bulk indexing, searching, and managing data streams.

- **Module:** `github.com/zmoog/es`
- **Go version:** 1.23 (toolchain 1.23.4)
- **Binary name:** `es`

## Repository Structure

```
es/
├── main.go                    # Entry point — calls cmd.Execute()
├── cmd/                       # CLI layer (Cobra commands, flag parsing)
│   ├── root.go                # Root command, Viper config init, global flags
│   ├── docs/
│   │   ├── docs.go            # 'es docs' subcommand group
│   │   ├── index.go           # 'es docs index' — index single document
│   │   └── bulk.go            # 'es docs bulk' — bulk index documents
│   ├── search/
│   │   └── search.go          # 'es search' — search documents
│   ├── datastream/
│   │   ├── datastream.go      # 'es datastream' subcommand group
│   │   └── delete.go          # 'es datastream delete' — delete data streams
│   └── version/
│       └── version.go         # 'es version' — print version info
├── es/                        # Core business logic package
│   ├── runner.go              # Runner: creates ES client and executes commands
│   └── commands/
│       ├── commands.go        # Command interface and UnitOfWork struct
│       ├── docs.go            # IndexDocCommand, BulkCommand implementations
│       ├── search.go          # SearchCommand implementation
│       ├── datastream.go      # DatastreamDeleteCommand implementation
│       └── version.go         # VersionCommand implementation
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yaml           # Release build configuration
├── .github/workflows/
│   └── release.yml            # CI/CD: triggered on version tags or manual dispatch
└── README.md
```

## Architecture

The project uses a **Command pattern** with a two-layer architecture:

1. **`cmd/` (CLI layer):** Cobra-based command definitions. Responsible for parsing flags, reading stdin/files, constructing command structs, creating a `Runner`, and calling `runner.Run(command)`.

2. **`es/` (core layer):** Business logic. The `Runner` creates the Elasticsearch client (from Viper config) and a `UnitOfWork`, then passes it to each `Command` via `ExecuteWith(uow)`.

### Key Types

```go
// Command interface — all commands implement this
type Command interface {
    ExecuteWith(uow UnitOfWork) error
}

// UnitOfWork — dependency container passed to commands
type UnitOfWork struct {
    Client *elasticsearch.Client
}

// Runner — orchestrates client setup and command execution
type Runner struct {
    uow commands.UnitOfWork
}
```

### Adding a New Command

1. Create `es/commands/<name>.go` with a struct implementing `Command` (i.e., `ExecuteWith(uow UnitOfWork) error`).
2. Create `cmd/<name>/` with a `NewCommand() *cobra.Command` function that wires flags, reads input, constructs the command struct, and calls `runner.Run(cmd)`.
3. Register the new Cobra command in `cmd/root.go` via `rootCmd.AddCommand(...)`.

## Development Workflows

### Run locally

```bash
# Run with a sample document
make run

# Or directly
go run main.go docs index --doc '{"@timestamp": "2024-01-01T00:00:00Z", "message": "hello"}'
```

### Build

```bash
go build -o es .
```

### Dependency management

```bash
make tidy
# equivalent to: go mod tidy
```

### Release

Releases are triggered automatically by pushing a git tag matching `v[0-9]+.[0-9]+.[0-9]+*`, or manually via GitHub Actions workflow dispatch. GoReleaser builds static binaries for Linux and macOS (amd64 and arm64) and publishes a Homebrew cask to `zmoog/homebrew-es`.

```bash
# Create a release tag
git tag v0.x.y
git push origin v0.x.y
```

Version information is injected at build time via ldflags into the `es/commands` package variables: `version`, `commit`, `date`, `builtBy`.

### Testing

There are currently no automated tests (`*_test.go` files). When adding tests, use:

```bash
go test ./...
```

## Configuration

Configuration is handled by Viper with the following priority (highest wins):

1. **CLI flags** — e.g., `--api-endpoints`, `--api-key`
2. **Environment variables** — prefix `ES_`, dots replaced by underscores (e.g., `ES_API_ENDPOINTS`, `ES_API_KEY`)
3. **Config file** — YAML at `~/.es/config` (or path specified via `--config`)

### Config file format (`~/.es/config`)

```yaml
api:
  endpoints: https://localhost:9200
  key: your-api-key

client:
  max-retries: 1
  retry-on-status:
    - 502
    - 503
    - 504
    - 429
  ca-cert-path: /path/to/ca-cert.pem  # optional, for self-signed certs
```

### Viper key mapping

| Viper key | CLI flag | Env var |
|---|---|---|
| `api.endpoints` | `--api-endpoints` / `-e` | `ES_API_ENDPOINTS` |
| `api.key` | `--api-key` / `-k` | `ES_API_KEY` |
| `client.max-retries` | `--client-max-retries` / `-m` | `ES_CLIENT_MAX_RETRIES` |
| `client.retry-on-status` | `--client-retry-on-status` / `-r` | `ES_CLIENT_RETRY_ON_STATUS` |
| `client.ca-cert-path` | `--client-ca-cert-path` / `-c` | `ES_CLIENT_CA_CERT_PATH` |

## Key Conventions

### Stdin input

Commands that accept a document, file, or query support reading from stdin by using `-` as the value:

```bash
cat docs.json | es docs bulk -f -
echo '{"query": {"match_all": {}}}' | es search my_index
```

In the CLI layer, check if the value is `"-"` (or empty) and read from `os.Stdin`. In the command struct, embed `io.Reader` (as in `BulkCommand`) or accept a `string` and handle the `-` sentinel.

### Error handling

Use `fmt.Errorf("context message: %w", err)` to wrap errors. The root command uses a `must()` helper for fatal Viper/flag binding errors.

### Retry strategy

The Elasticsearch client uses exponential backoff via `github.com/cenkalti/backoff/v4`. Retry status codes default to 502, 503, 504, 429 and are configurable.

### Bulk indexing

Use `esutil.NewBulkIndexer` from `github.com/elastic/go-elasticsearch/v8/esutil`. Read input line-by-line (each line = one JSON document) using `bufio.Scanner`. Provide `OnSuccess` and `OnFailure` callbacks and log stats after `bulkIndexer.Close(ctx)`.

### Search output

Search results are written directly to `os.Stdout` via `io.Copy`. Pipe through `jq` for formatting.

## Key Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI command framework |
| `github.com/spf13/viper` | Configuration management |
| `github.com/elastic/go-elasticsearch/v8` | Elasticsearch Go client |
| `github.com/cenkalti/backoff/v4` | Exponential backoff for retries |
| `github.com/zmoog/classeviva` | Error/feedback logging (`feedback.Error`) |
