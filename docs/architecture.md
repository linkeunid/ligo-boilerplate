# Architecture

## Overview

This boilerplate follows a modular architecture using the Ligo framework. Each feature is organized as an independent module with its own controllers, services, and repositories.

## Layer Structure

```
┌─────────────────────────────────────────────────────┐
│                    Controllers                       │
│  (HTTP handlers, request/response, route guards)    │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│                     Services                         │
│           (Business logic, validation)               │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│                   Repositories                       │
│              (Data access, storage)                  │
└─────────────────────────────────────────────────────┘
```

## Request Flow

```
Request → Middleware → Guards → Interceptors → Controller → Service → Repository
                                    ↑              ↓
                                  Filters         Response
```

### Middleware

Global middleware applies to all routes:
- `RecoveryMiddleware` - Panic recovery
- `CORSMiddleware` - Cross-origin headers

### Guards

Authorization checks before handler execution:
- `AuthGuard` - Requires valid JWT
- `AdminGuard` - Requires admin role
- `RolesGuard` - Custom role-based guard

### Interceptors

Around-advice for cross-cutting concerns:
- `LoggingInterceptor` - Request/response logging
- `AuditInterceptor` - Admin action audit trail

### Filters

Exception handling and response transformation:
- `GlobalExceptionFilter` - Converts errors to HTTP responses

## Dependency Injection

Ligo's container manages dependencies automatically:

```go
ligo.Providers(
    ligo.Factory[*UserService](NewUserService),
    ligo.Factory[*UserRepository](NewUserRepository),
)
```

Controllers receive dependencies via constructor injection.

## Module System

Each module exports:
- **Providers** - Services, repositories, dependencies
- **Controllers** - HTTP handlers
- **Imports** - Required modules

```go
func Module() ligo.Module {
    return ligo.NewModule("user",
        ligo.Imports(auth.Module()),
        ligo.Providers(...),
        ligo.Controllers(...),
    )
}
```

## Common Package

Shared utilities in `internal/common/`:

| File | Purpose |
|------|---------|
| `errors.go` | Standard error types |
| `filter.go` | Global exception filter |
| `interceptor.go` | Logging interceptor |
| `version.go` | Application version |
