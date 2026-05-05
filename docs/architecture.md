# Architecture

## Overview

This boilerplate follows **Clean Architecture + Repository Pattern** using the Ligo framework.
Business logic is isolated from infrastructure concerns — the domain layer has zero external dependencies.

## Layer Structure

```
┌──────────────────────────────────────────────────────────┐
│                        cmd/api/                           │
│           Entry point — wires all layers together         │
└──────────────────────────────────────────────────────────┘
                             ↓
┌──────────────────────────────────────────────────────────┐
│                     infrastructure/                       │
│  HTTP controllers, middleware, auth, persistence impls    │
└──────────────────────────────────────────────────────────┘
                             ↓
┌──────────────────────────────────────────────────────────┐
│                       usecase/                            │
│        Application business logic, orchestration          │
└──────────────────────────────────────────────────────────┘
                             ↓
┌──────────────────────────────────────────────────────────┐
│                        domain/                            │
│   Entities, repository interfaces, service interfaces     │
│              NO external dependencies                     │
└──────────────────────────────────────────────────────────┘
```

## Dependency Flow

```
cmd → infrastructure → usecase → domain
                                  ↑
                           (never reversed)
```

## Layer Responsibilities

| Layer | Package | Purpose | Depends On |
|-------|---------|---------|------------|
| Domain | `internal/domain/` | Entities, interfaces | Nothing |
| UseCase | `internal/usecase/` | Business logic, orchestration | domain |
| Infrastructure | `internal/infrastructure/` | HTTP, DB, auth, external services | domain, usecase |
| Entry | `cmd/api/` | Wiring, configuration | All |

## Directory Layout

```
internal/
├── domain/
│   ├── entity/        # Pure business entities (User, File)
│   ├── repository/    # Repository interfaces — contracts for data access
│   └── service/       # Domain service interfaces (AuthService, roles)
│
├── usecase/
│   ├── user.go        # User business logic
│   ├── file.go        # File business logic
│   ├── errors.go      # Common use case errors
│   └── dto/
│       ├── create_user.go
│       └── update_user.go
│
├── infrastructure/
│   ├── auth/
│   │   ├── jwt.go     # JWTAuth — implements domain/service.AuthService
│   │   └── guard.go   # AuthGuard, AdminGuard
│   ├── http/
│   │   ├── controller/   # HTTP handlers — call usecases, use pipes for validation
│   │   ├── middleware/   # Exception handling, logging, audit
│   │   └── presenter/    # Response formatting
│   └── persistence/
│       └── memory/    # In-memory repository implementations
│           ├── user_repo.go  # implements domain/repository.UserRepository
│           └── file_repo.go  # implements domain/repository.FileRepository
│
├── module/            # Wires layers into Ligo modules
└── config/            # Application configuration
```

## Request Flow

```
HTTP Request
    → Global Middleware (CORS, recovery)
    → Route Middleware (logging, exception)
    → Guards (auth, admin)
    → Pipes (UUID validation, body binding + struct validation)
    → Controller (retrieve validated data, call use case)
    → UseCase (business logic)
    → Repository (data access)
    → Response
```

## Middleware Behaviour

`LoggingMiddleware` logs at different levels depending on the error type:

- 4xx errors (client errors) — logged as `WARN`
- 5xx errors (server errors) — logged as `ERROR`
- No error — logged as `DEBUG`

`ExceptionMiddleware` catches all errors, logs them, and maps them to HTTP responses. See [Authentication](authentication.md) for the full error → status code mapping.

## Lifecycle Hooks

The boilerplate demonstrates Ligo's lifecycle hook system at three levels:

- **Application-level** (`OnStart`/`OnStop` in `main.go`)
- **Module-level** (`OnModuleInit`/`OnModuleDestroy` in module definitions)
- **Provider/Controller-level** (via interface implementations or `HookedFactory`/`HookedController`)

See [Lifecycle Hooks](lifecycle.md) for details on hook patterns, execution order, and examples.

## Key Principles

**1. Domain isolation** — `internal/domain/` imports nothing from this project. Entities are plain structs; repository contracts are interfaces.

**2. Dependency inversion** — UseCases depend on repository *interfaces* (defined in domain), not concrete implementations. Swap `memory` → `postgres` without touching business logic.

**3. Factory functions for controllers** — Ligo's `ligo.Controllers()` accepts constructors directly. The DI binder resolves parameters by type from the container and calls the constructor via reflection:

```go
ligo.Controllers(controller.NewUserController)
```

**4. Providers wire the dependency graph** — Each module registers its providers in `ligo.Providers(...)`. The DI container resolves the full graph at startup:

```go
ligo.Providers(
    ligomemory.Provider[string, *entity.User](),
    ligo.Factory[repository.UserRepository](memory.NewUserRepository),
    ligo.Factory[*usecase.UserUseCase](usecase.NewUserUseCase),
),
ligo.Controllers(controller.NewUserController),
```

**5. Validation is a pipe concern** — Use case methods receive already-validated DTOs via `ligo.ValidatedBody[T]`. Use cases do not instantiate a validator or re-validate inputs.
