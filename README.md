# Ligo Boilerplate

> Go web application boilerplate using [Ligo](https://github.com/linkeunid/ligo) framework.

## Quick Start

```bash
# Run
go run cmd/example/main.go

# Build
go build -o bin/app cmd/example/main.go

# Run binary
./bin/app
```

Server runs on `http://localhost:8080`.

## API Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET | `/` | API info | - |
| GET | `/health` | Health check | - |
| GET | `/users` | List users | - |
| GET | `/users/:id` | Get user | Bearer token |
| POST | `/users` | Create user | Bearer token |
| PUT | `/users/:id` | Update user | Bearer token |
| DELETE | `/users/:id` | Delete user | Admin only |
| POST | `/files/upload` | Upload file | - |
| GET | `/files/:id` | Download file | - |
| GET | `/files` | List files | - |
| DELETE | `/files/:id` | Delete file | - |

## Example Auth Tokens

```
user:secret   (regular user)
admin:secret  (admin user)
```

```bash
curl -H "Authorization: Bearer user:secret" http://localhost:8080/users/1
```

## Project Structure

```
internal/
├── auth/      # Guards, interceptors
├── common/    # Shared utilities
├── file/      # File upload
├── health/    # Health check
├── root/      # API info
└── user/      # User CRUD
```

## Documentation

- [Architecture](docs/architecture.md)
- [Modules Guide](docs/modules.md)
- [Authentication](docs/authentication.md)
- [Contributing](docs/contributing.md)

## License

MIT
