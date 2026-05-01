package todo

import (
	"github.com/linkeunid/ligo"
)

// Controller handles HTTP requests for todos
type Controller struct {
	svc *Service
}

// NewController creates a new todo controller
func NewController(svc *Service) ligo.Controller {
	return &Controller{svc: svc}
}

// Routes registers all routes for the todo module
func (c *Controller) Routes(r ligo.Router) {
	api := r.Group("/api/todos")

	api.Handle("GET", "", c.List)
	api.Handle("POST", "", c.Create)
	api.Handle("GET", "/:id", c.Get)
	api.Handle("PUT", "/:id", c.Update)
	api.Handle("DELETE", "/:id", c.Delete)
	api.Handle("PATCH", "/:id/toggle", c.Toggle)
}

// List handles GET /api/todos
func (c *Controller) List(ctx ligo.Context) error {
	todos := c.svc.List()
	return ctx.JSON(200, map[string]any{
		"data": todos,
		"count": len(todos),
	})
}

// Create handles POST /api/todos
func (c *Controller) Create(ctx ligo.Context) error {
	var req CreateTodoRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request body"})
	}

	todo, err := c.svc.Create(req)
	if err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(201, map[string]any{"data": todo})
}

// Get handles GET /api/todos/:id
func (c *Controller) Get(ctx ligo.Context) error {
	id := ctx.Param("id")

	todo, err := c.svc.GetByID(id)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{"data": todo})
}

// Update handles PUT /api/todos/:id
func (c *Controller) Update(ctx ligo.Context) error {
	id := ctx.Param("id")

	var req UpdateTodoRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request body"})
	}

	todo, err := c.svc.Update(id, req)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{"data": todo})
}

// Delete handles DELETE /api/todos/:id
func (c *Controller) Delete(ctx ligo.Context) error {
	id := ctx.Param("id")

	if err := c.svc.Delete(id); err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]string{"message": "deleted"})
}

// Toggle handles PATCH /api/todos/:id/toggle
func (c *Controller) Toggle(ctx ligo.Context) error {
	id := ctx.Param("id")

	todo, err := c.svc.Toggle(id)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{"data": todo})
}
