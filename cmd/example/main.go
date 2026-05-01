package main

import (
	"net/http"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo/adapters/echo"
	"github.com/linkeunid/ligo-boilerplate/internal/health"
	"github.com/linkeunid/ligo-boilerplate/internal/root"
	"github.com/linkeunid/ligo-boilerplate/internal/user"
)

// CORS middleware for cross-origin requests.
func CORSMiddleware(next ligo.HandlerFunc) ligo.HandlerFunc {
	return func(ctx ligo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request().Method == "OPTIONS" {
			return ctx.String(204, "")
		}

		return next(ctx)
	}
}

// Recovery middleware for panic recovery.
func RecoveryMiddleware(next ligo.HandlerFunc) ligo.HandlerFunc {
	return func(ctx ligo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				ctx.Response().WriteHeader(500)
			}
		}()
		return next(ctx)
	}
}

func main() {
	// Create Echo adapter
	router := echo.NewAdapter()

	// Create Ligo application
	app := ligo.New(
		ligo.WithRouter(router),
		ligo.WithAddr(":8080"),
		ligo.WithAutoPort(), // Enable automatic port increment if port is in use
		// Global middleware - applies to all routes
		ligo.WithMiddleware(
			RecoveryMiddleware,
			CORSMiddleware,
		),
		// Lifecycle hooks
		ligo.OnStart(func(ctx any) error {
			return nil
		}),
		ligo.OnStop(func(ctx any) error {
			return nil
		}),
	)

	// Register modules
	// Auth module is imported by user.Module(), no need to register separately
	app.Register(
		user.Module(),  // User module (imports auth.Module())
		health.Module(), // Health check endpoints
		root.Module(),   // API info endpoint
	)

	// Start the server
	if err := app.Run(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
