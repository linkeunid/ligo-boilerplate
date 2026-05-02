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
│   │   ├── controller/   # HTTP handlers — call usecases
│   │   ├── middleware/   # Exception handling, logging, audit
│   │   ├── presenter/    # Response formatting
│   │   └── validator/    # Request validation
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
    → Controller (parse request)
    → UseCase (business logic)
    → Repository (data access)
    → Response
```

## Key Principles

**1. Domain isolation** — `internal/domain/` imports nothing from this project. Entities are plain structs; repository contracts are interfaces.

**2. Dependency inversion** — UseCases depend on repository *interfaces* (defined in domain), not concrete implementations. Swap `memory` → `postgres` without touching business logic.

**3. Factory functions for controllers** — Ligo's `ligo.Controllers()` expects factory functions so the DI binder can call them via reflection:

```go
ligo.Controllers(
    func() ligo.Controller { return controller.NewUserController(uc, log) },
)
```

**4. Manual wiring in modules** — Each module function (`module.User(cfg, log)`) manually constructs its dependency graph. This is explicit and testable without a DI framework.
