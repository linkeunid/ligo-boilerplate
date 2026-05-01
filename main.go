package main

import (
	"fmt"
	"net/http"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo/adapters/echo"
	"github.com/linkeunid/ligo-example/internal/middleware"
	"github.com/linkeunid/ligo-example/internal/todo"
)

func main() {
	router := echo.NewAdapter()
	app := ligo.New(
		ligo.WithRouter(router),
		ligo.WithAddr(":8080"),
		ligo.WithMiddleware(
			middleware.CORS,
			middleware.RequestID,
		),
		ligo.WithDebug(true),
	)

	// Register modules - config module imported by todo module via ligo.Imports()
	// so only todo needs to be registered here
	app.Register(todo.Module())

	fmt.Println("========================================")
	fmt.Println("  Ligo Example - Todo API")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /api/todos          - List all todos")
	fmt.Println("  GET    /api/todos/:id      - Get todo by ID")
	fmt.Println("  POST   /api/todos          - Create todo")
	fmt.Println("  PUT    /api/todos/:id      - Update todo")
	fmt.Println("  DELETE /api/todos/:id      - Delete todo")
	fmt.Println("  PATCH  /api/todos/:id/toggle - Toggle completion")
	fmt.Println()
	fmt.Println("Try:")
	fmt.Println(`  curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{"title":"Build something with ligo"}'`)
	fmt.Println(`  curl http://localhost:8080/api/todos`)
	fmt.Println()

	if err := app.Run(); err != nil {
		if err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}
}
