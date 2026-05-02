# Contributing

## Setup

```bash
# Clone repository
git clone https://github.com/linkeunid/ligo-boilerplate.git
cd ligo-boilerplate

# Install dependencies
go mod download

# Run
go run ./cmd/api/
```

## Development

### Project Structure

```
.
├── cmd/api/                    # Application entry point
├── internal/                   # Private application code
│   ├── config/                 # Application configuration
│   ├── domain/                 # Core business (no external deps)
│   │   ├── entity/             # Business entities
│   │   ├── repository/         # Repository interfaces
│   │   └── service/            # Domain service interfaces
│   ├── usecase/                # Application business logic
│   │   └── dto/                # Data Transfer Objects
│   ├── infrastructure/         # External concerns
│   │   ├── auth/               # JWT implementation + guards
│   │   ├── http/               # Controllers, middleware, presenters
│   │   └── persistence/        # Repository implementations
│   └── module/                 # Module wiring (connects layers)
├── docs/                       # Documentation
└── go.mod                      # Go module definition
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Validation belongs at the HTTP layer (pipes), not in use cases

### Adding a Feature

1. Define entity and repository interface in `internal/domain/`
2. Implement business logic in `internal/usecase/`
3. Implement repository in `internal/infrastructure/persistence/memory/`
4. Implement controller in `internal/infrastructure/http/controller/`
5. Wire the module in `internal/module/`
6. Register in `cmd/api/main.go`
7. Add tests
8. Update documentation

See [Modules Guide](modules.md) for a step-by-step walkthrough.

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/usecase/...
```

### Commit Messages

```
feat: add product module
fix: handle missing user in update
docs: update authentication guide
refactor: move auth guard to infrastructure/auth
```

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a PR
