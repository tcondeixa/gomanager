# CRUSH Configuration - goinstall

## Build/Test/Lint Commands
- **Build**: `make build.local` or `go build .`
- **Test**: `make test` (runs `go test -v -race -cover ./...`)
- **Test single package**: `go test -v -race -cover ./internal/pkg`
- **Lint**: `make lint` (golangci-lint with fix) or `make lint.ci` (CI mode)
- **Format**: `make fmt` (gofmt + golines with 120 char limit)
- **Pre-commit**: `make pre-commit` (fmt + lint)
- **Clean**: `make clean`

## Code Style Guidelines
- **Go version**: 1.25.1+
- **Line length**: 120 characters (enforced by golines)
- **Formatting**: gofmt + gofumpt + goimports (via golangci-lint)
- **Imports**: Standard library first, then third-party, then local packages
- **Naming**: camelCase for private, PascalCase for public, descriptive names
- **Error handling**: Always check errors, wrap with context using fmt.Errorf
- **Logging**: Use log/slog for structured logging, not fmt.Print*
- **Constants**: Use const blocks, ALL_CAPS for package-level constants
- **File permissions**: Use octal notation (0o755, 0o644)
- **JSON tags**: Always include json tags for structs that are serialized
- **Generics**: Use meaningful type parameter names, prefer [T any] for simple cases

## Project Structure
- `cmd/`: Cobra CLI commands
- `internal/pkg/`: Package management logic  
- `internal/storage/`: Generic storage layer with JSON persistence
- `main.go`: Entry point with version string

## Linting Rules
- errcheck, govet, ineffassign, revive, staticcheck, unused enabled
- govet: fieldalignment and shadow disabled
- Common false-positives and legacy issues excluded