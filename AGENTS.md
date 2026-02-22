# AGENTS.md - Agentic Coding Guidelines for a_c_a_c

## Project Overview

a_c_a_c is a Twitch chatbot that generates responses using a Markov chain algorithm. It learns from chat messages in real-time and generates responses. The project is written in Go (1.25+).

## Build, Lint, and Test Commands

### Building

```bash
# Build the binary
go build -o acac .

# Build for Docker
docker build -t ghcr.io/notarock/a_c_a_c:latest .
```

### Running

```bash
# Run with environment variables (see .env.example)
./acac

# Run single message generation from file
./acac -from-file /path/to/messages.txt
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./pkg/chain

# Run tests matching a pattern
go test -v -run TestFunctionName ./...

# Run tests with coverage
go test -v -cover ./...
```

### Linting

```bash
# Run go vet
go vet ./...

# Run golint (if installed)
golangci-lint run ./...

# Format code
go fmt ./...

# All-in-one (fmt + vet)
go fmt ./... && go vet ./...
```

## Code Style Guidelines

### General Principles

- Write clean, readable code with clear variable and function names
- Keep functions small and focused on a single responsibility
- Use early returns to reduce nesting
- Handle errors explicitly - avoid silently ignoring errors

### Naming Conventions

- **Packages**: lowercase, short, descriptive (e.g., `chain`, `filters`, `twitch`)
- **Types/Structs**: PascalCase (e.g., `Chain`, `TwitchClient`, `Filter`)
- **Functions/Methods**: PascalCase (e.g., `NewChain`, `LoadModel`)
- **Variables/Fields**: camelCase (e.g., `savedMessagesFilepath`, `sending`)
- **Constants**: PascalCase for exported, can use camelCase or UPPER_SNAKE_CASE for unexported (e.g., `GREEN`, `RED`)
- **Interfaces**: PascalCase, typically with "er" suffix (e.g., `Filter`, `Runner`)

### Imports

- Group imports: standard library first, then third-party packages
- Use blank import (`_`) only when necessary (e.g., `godotenv/autoload`)
- Sort imports alphabetically within groups

```go
import (
    "flag"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/joho/godotenv/autoload"
    "github.com/notarock/a_c_a_c/pkg/chain"
)
```

### Error Handling

- Use `fmt.Errorf` with `%v` for wrapping errors with context
- Return errors rather than using `log.Panic` unless it's a fatal startup error
- Check errors at call site and handle appropriately

```go
// Good
if err != nil {
    return fmt.Errorf("failed to load previous messages from file %s: %v", filepath, err)
}

// For fatal startup errors
if BASE_PATH == "" {
    log.Panic("Missing environment variables")
}
```

### Structs and Types

- Use struct tags for YAML serialization (use `yaml:` prefix)
- Define config structs separately from business logic structs
- Use pointer receivers only when modifying the struct

```go
type ChainConfig struct {
    SentMessagesFilepath  string
    SavedMessagesFilepath string
    Saving                bool
    IgnoreParrots         bool
}

type Channel struct {
    Name      string   `yaml:"name"`
    Frequency int      `yaml:"frequency"`
    AllowBits bool     `yaml:"allow_bits"`
}
```

### Interfaces

- Define small, focused interfaces
- Accept interfaces, return concrete types when possible
- Keep interface definitions in the same package that uses them

```go
type Filter interface {
    Filter(message string) bool
}
```

### Logging

- Use `log` for errors and warnings
- Use `fmt.Println` for informational output during startup/development
- Use environment checks (`ENV != "production"`) for debug output

### Concurrency

- Use goroutines with channels for concurrent operations
- Be careful with shared state - use mutexes or channels as appropriate

### Comments

- Document exported functions and types with doc comments
- Keep comments concise and descriptive
- Use Go-style doc comments (start with the name of the element)

### Code Organization

- Main logic in `main.go`
- Package code in `pkg/<name>/` directories
- One file per type/feature when reasonable
- Group related functionality (e.g., filter implementations in `filters/`)

### Configuration

- Use environment variables for runtime configuration
- Use YAML files for channel-specific configuration
- Follow the existing pattern in `config.LoadChannelConfig`

## Testing Guidelines

- Write tests for new functionality (no tests currently exist)
- Use table-driven tests when testing multiple cases
- Test both success and error paths
- Place tests in `*_test.go` files in the same package

## CI/CD

The project uses "tuyauterie" workflow (defined in `.tuyauterie.yaml`):
- Builds Go binary and Docker image
- Runs on push, pull request, release, and schedule
- Uses GitOps for deployment

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Style Guide](https://google.github.io/styleguide/go/)
