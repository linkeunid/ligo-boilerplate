# Modules Guide

## Creating a New Module

Follow Clean Architecture: define the domain first, then work outward.

### 1. Domain — entity and repository interface

```go
// internal/domain/entity/product.go
package entity

type Product struct {
    ID    string
    Name  string
    Price float64
}
```

```go
// internal/domain/repository/product.go
package repository

import "github.com/linkeunid/ligo-boilerplate/internal/domain/entity"

type ProductRepository interface {
    FindByID(id string) (*entity.Product, bool)
    FindAll() []*entity.Product
    Create(name string, price float64) *entity.Product
    Delete(id string) bool
}
```

### 2. UseCase — business logic

```go
// internal/usecase/product.go
package usecase

import (
    "github.com/linkeunid/ligo"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

type ProductUseCase struct {
    repo repository.ProductRepository
    log  ligo.Logger
}

func NewProductUseCase(repo repository.ProductRepository, log ligo.Logger) *ProductUseCase {
    return &ProductUseCase{repo: repo, log: log}
}

func (uc *ProductUseCase) GetProduct(id string) (*entity.Product, error) {
    p, found := uc.repo.FindByID(id)
    if !found {
        return nil, ErrNotFound
    }
    return p, nil
}
```

### 3. Persistence — repository implementation

```go
// internal/infrastructure/persistence/memory/product_repo.go
package memory

import (
    "sync"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

type ProductRepository struct {
    mu       sync.RWMutex
    products map[string]*entity.Product
}

func NewProductRepository() repository.ProductRepository {
    return &ProductRepository{products: make(map[string]*entity.Product)}
}

func (r *ProductRepository) FindByID(id string) (*entity.Product, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    p, found := r.products[id]
    return p, found
}

// ... implement remaining interface methods
```

### 4. Controller — HTTP handler

```go
// internal/infrastructure/http/controller/product.go
package controller

import (
    "github.com/linkeunid/ligo"
    "github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

type ProductController struct {
    uc  *usecase.ProductUseCase
    log ligo.Logger
}

func NewProductController(uc *usecase.ProductUseCase, log ligo.Logger) *ProductController {
    return &ProductController{uc: uc, log: log}
}

func (c *ProductController) Routes(r ligo.Router) {
    cr := ligo.NewChainRouter(r.Group("/products"))
    cr.GET("", c.ListProducts).Handle()
    cr.GET("/:id", c.GetProduct).Handle()
}

func (c *ProductController) GetProduct(ctx ligo.Context) error {
    id := ctx.Param("id")
    product, err := c.uc.GetProduct(id)
    if err != nil {
        return err
    }
    return ctx.OK(product)
}
```

### 5. Module — wire the layers

```go
// internal/module/product.go
package module

import (
    "github.com/linkeunid/ligo"
    "github.com/linkeunid/ligo-boilerplate/internal/config"
    "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
    "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/middleware"
    "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/persistence/memory"
    "github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

func Product(cfg *config.Config, log ligo.Logger) ligo.Module {
    repo := memory.NewProductRepository()
    uc := usecase.NewProductUseCase(repo, log)

    return ligo.NewModule("product",
        ligo.Controllers(
            func() ligo.Controller { return controller.NewProductController(uc, log) },
        ),
    )
}
```

> **Note:** `ligo.Controllers()` requires factory functions — wrap constructors in a closure. The binder calls the function via reflection to bind routes.

### 6. Register in main

```go
// cmd/api/main.go
app.Register(
    module.Product(cfg, log),
    // ...
)
```

## Adding Auth Guards

```go
import (
    infraauth "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/auth"
    "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/middleware"
)

func Product(cfg *config.Config, log ligo.Logger) ligo.Module {
    repo := memory.NewProductRepository()
    uc := usecase.NewProductUseCase(repo, log)
    authSvc := infraauth.NewJWTAuth(log)

    type productRoutes struct {
        ctrl        *controller.ProductController
        authGuard   ligo.Guard
        exceptionMW ligo.Middleware
    }

    routes := &productRoutes{
        ctrl:        controller.NewProductController(uc, log),
        authGuard:   infraauth.AuthGuard(authSvc),
        exceptionMW: middleware.ExceptionMiddleware(log),
    }

    return ligo.NewModule("product",
        ligo.Controllers(
            func() ligo.Controller {
                return &struct{ *productRoutes }{routes}
            },
        ),
    )
}
```

## Existing Modules

| Module | Function | Endpoints |
|--------|---------|-----------|
| `auth` | `module.Auth(cfg, log)` | Provides `AuthService` |
| `user` | `module.User(cfg, log)` | `/users/*` |
| `file` | `module.File(cfg, log)` | `/files/*` |
| `health` | `module.Health(cfg)` | `/health` |
| `root` | `module.Root(cfg)` | `/` |

## Swapping Persistence

The repository pattern means replacing the storage backend requires only changing the concrete implementation in the module — the domain and usecase layers are untouched:

```go
// From in-memory:
repo := memory.NewProductRepository()

// To PostgreSQL (when implemented):
repo := postgres.NewProductRepository(db)
```
