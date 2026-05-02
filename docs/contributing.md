# Contributing

## Setup

```bash
# Clone repository
git clone https://github.com/linkeunid/ligo-boilerplate.git
cd ligo-boilerplate

# Install dependencies
go mod download

# Run
go run cmd/example/main.go
```

## Development

### Project Structure

```
.
├── cmd/example/       # Application entry point
├── internal/          # Private application code
│   ├── auth/         # Authentication
│   ├── common/       # Shared utilities
│   ├── file/         # File handling
│   ├── health/       # Health checks
│   ├── root/         # Root endpoint
│   └── user/         # User module
├── docs/             # Documentation
└── go.mod            # Go module definition
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Export types via `ligo.Factory` for DI

### Adding a Feature

1. Create module in `internal/feature/`
2. Implement repository, service, controller
3. Define module with providers/controllers
4. Register in `cmd/example/main.go`
5. Add tests
6. Update documentation

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/user
```

### Commit Messages

```
feat: add file upload module
fix: resolve import cycle in common package
docs: update authentication guide
refactor: move interceptors to common package
```

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a PR

## Version Management

Update `internal/common/version.go` for release bumps:

```go
const Version = "0.8.0"
```
