# Lifecycle Hooks

Ligo provides a comprehensive lifecycle hook system with 5 stages for managing application initialization and shutdown.

**See Also:** [Before Shutdown Example](examples/before_shutdown_example.go)

## Hook Stages

| Stage | When Called | Use Case |
|-------|-------------|----------|
| `OnModuleInit` | Module initializes (depth-first) | Resource initialization |
| `OnApplicationBootstrap` | After all modules init, before serving | Readiness checks |
| `BeforeApplicationShutdown` | Before shutdown begins | Drain-stop (stop accepting new requests) |
| `OnApplicationShutdown` | During shutdown | Cleanup resources |
| `OnModuleDestroy` | Module destroys (reverse depth-first) | Final cleanup |

## Hook Patterns

### 1. Interface-Based (Duck-Typing)

Implement hook interfaces directly on providers/controllers:

```go
type JWTAuth struct {
    log ligo.Logger
}

func (j *JWTAuth) OnModuleInit() error {
    j.log.Info("JWT authentication initialized")
    return nil
}

func (j *JWTAuth) OnApplicationShutdown() error {
    j.log.Info("JWT authentication shutting down")
    return nil
}
```

**Used in:** `internal/infrastructure/auth/jwt.go`, `internal/infrastructure/persistence/memory/file_repo.go`

### 2. Compile-Time Safe (HookedFactory/HookedController)

Use the `Registerable` interface for compile-time safe hook registration:

```go
type UserRepository struct {
    store *ligomemory.Store[int, *entity.User]
    log  ligo.Logger
}

func (r *UserRepository) Register(registry *ligo.HookRegistry) {
    registry.OnInit(r.SeedDatabase)
    registry.OnDestroy(r.CleanupDatabase)
}

func (r *UserRepository) SeedDatabase() error {
    r.log.Info("Seeding user database")
    return nil
}

func (r *UserRepository) CleanupDatabase() error {
    r.log.Info("Cleaning up user database")
    return nil
}
```

**Register with HookedFactory:**
```go
ligo.Providers(
    ligo.HookedFactory[repository.UserRepository](memory.NewUserRepository),
)
```

**Used in:** `internal/infrastructure/persistence/memory/user_repo.go`, all controllers

### 3. Module-Level Hooks

Register hooks when defining a module:

```go
func FileModule() ligo.Module {
    return ligo.NewModule("file",
        ligo.Providers(...),
        ligo.OnModuleInit(func() error {
            // Module-level initialization
            return nil
        }),
        ligo.OnModuleDestroy(func() error {
            // Module-level cleanup
            return nil
        }),
    )
}
```

**Note:** Module-level hooks don't support DI injection. Use provider-level hooks for logging.

## Controller Hooks

Controllers implement hooks via the `Registerable` interface with `HookedController`:

```go
type FileController struct {
    uc  *usecase.FileUseCase
    log ligo.Logger
}

func (c *FileController) Initialize() error {
    c.log.Info("File controller initializing")
    return nil
}

func (c *FileController) Ready() error {
    c.log.Info("File controller ready")
    return nil
}

func (c *FileController) Draining() error {
    c.log.Info("File controller draining")
    return nil
}

func (c *FileController) Shutdown() error {
    c.log.Info("File controller shutting down")
    return nil
}

func (c *FileController) Register(registry *ligo.HookRegistry) {
    registry.OnInit(c.Initialize)
    registry.OnBootstrap(c.Ready)
    registry.BeforeShutdown(c.Draining)
    registry.OnShutdown(c.Shutdown)
}
```

**Register with HookedController:**
```go
ligo.Controllers(ligo.HookedController(controller.NewFileController))
```

**Benefits:**
- ✅ Compile-time safety for method names
- ✅ Descriptive method names (`Initialize` vs `OnModuleInit`)
- ✅ Explicit registration shows intent

## Application-Level Hooks

Registered in `cmd/api/main.go`:

```go
app := ligo.New(
    ligo.OnStart(func(ctx any) error {
        log.Info("Server starting")
        return nil
    }),
    ligo.OnStop(func(ctx any) error {
        log.Info("Server stopped gracefully")
        return nil
    }),
)
```

## Execution Order

### Startup
```
1. Application OnStart
2. Module OnModuleInit (depth-first)
3. Provider OnModuleInit
4. Controller OnModuleInit
5. Provider OnApplicationBootstrap
6. Controller OnApplicationBootstrap
7. Server starts serving
```

### Shutdown
```
1. Shutdown signal received
2. Server stops accepting connections
3. Provider BeforeApplicationShutdown (reverse order)
4. Controller BeforeApplicationShutdown (reverse order)
5. Provider OnApplicationShutdown (reverse order)
6. Controller OnApplicationShutdown (reverse order)
7. Provider OnModuleDestroy (reverse order)
8. Module OnModuleDestroy (reverse order)
9. Application OnStop
```

## Pattern Comparison

| Pattern | Compile-Time Safe | DI Support | Best For |
|---------|-------------------|------------|----------|
| Interface-Based | ❌ | ✅ | Simple providers |
| HookedFactory/HookedController | ✅ | ✅ | Production services |
| Module-Level | ✅ | ❌ | Module-wide coordination |

## Testing

```bash
go run cmd/api/main.go
```

**Startup output:**
```
INFO Server starting addr=:8080
INFO JWT authentication initialized
INFO Seeding user database with initial data
INFO File controller initializing upload_dir=./uploads max_file_size=10485760
INFO User controller initializing
INFO Health controller initializing
INFO Root controller initializing
INFO File controller ready to handle requests
INFO User controller ready to handle requests
INFO Health controller ready
INFO Root controller ready
INFO HTTP server started
```

**Shutdown output (Ctrl+C):**
```
INFO Root controller draining
INFO File controller draining - stopping new file uploads
INFO User controller draining - completing in-flight requests
INFO Root controller shutting down
INFO File controller shutting down
INFO User controller shutting down
INFO Cleaning up user database
INFO JWT authentication shutting down
INFO Server stopped gracefully
```

## Best Practices

1. **Use HookedFactory/HookedController** for production code (compile-time safety)
2. **Keep hooks simple** — fast and non-blocking
3. **Handle errors** — return errors to prevent startup/shutdown
4. **Use structured logging** — log with `ligo.LoggerField` for clarity
5. **Idempotent operations** — hooks may be called multiple times
6. **Timeout awareness** — don't block indefinitely

## Further Reading

- [Ligo Framework Documentation](https://github.com/linkeunid/ligo)
- [Ligo App & Lifecycle](../ligo/docs/features/app.md)
- [Ligo Best Practices](../ligo/docs/best-practices.md)
