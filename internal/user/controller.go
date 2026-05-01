package user

import (
	"errors"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/auth"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

// Controller handles HTTP routes for users.
type Controller struct {
	svc  *UserService
	log  ligo.Logger
	auth *auth.AuthService
}

// NewController creates a new user controller.
func NewController(svc *UserService, log ligo.Logger, auth *auth.AuthService) ligo.Controller {
	return &Controller{svc: svc, log: log, auth: auth}
}

// Routes registers all routes for the user module.
// Demonstrates all Ligo features: Guards, Pipes, Interceptors, Exception Filters.
func (c *Controller) Routes(r ligo.Router) {
	cr := ligo.NewChainRouter(r.Group("/users"))

	// Public routes - no authentication required
	cr.GET("", c.GetAllUsers).
		Intercept(LoggingInterceptor(c.log)).
		Handle()

	// Authenticated routes - require valid JWT
	cr.GET("/:id", c.GetUserByID).
		Guard(auth.AuthGuard(c.auth)).
		Intercept(LoggingInterceptor(c.log)).
		Handle()

	// Admin routes - require admin role
	cr.DELETE("/:id", c.DeleteUser).
		Guard(auth.AuthGuard(c.auth), auth.AdminGuard()).
		Intercept(LoggingInterceptor(c.log)).
		Handle()

	// Create user - requires auth + validation
	cr.POST("", c.CreateUser).
		Guard(auth.AuthGuard(c.auth)).
		Pipe(CreateUserValidationPipe(c.log)).
		Intercept(LoggingInterceptor(c.log)).
		Handle()

	// Update user - requires auth + validation
	cr.PUT("/:id", c.UpdateUser).
		Guard(auth.AuthGuard(c.auth)).
		Pipe(UpdateUserValidationPipe(c.log)).
		Intercept(LoggingInterceptor(c.log)).
		Handle()
}

// GetAllUsers handles GET /users
func (c *Controller) GetAllUsers(ctx ligo.Context) error {
	users := c.svc.GetAllUsers()

	return ctx.JSON(200, map[string]any{
		"users": users,
		"count": len(users),
	})
}

// GetUserByID handles GET /users/:id
func (c *Controller) GetUserByID(ctx ligo.Context) error {
	id := ctx.Param("id")

	user, err := c.svc.GetUserByID(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.JSON(404, map[string]string{"error": "user not found"})
		}
		return err
	}

	return ctx.JSON(200, user)
}

// CreateUser handles POST /users
func (c *Controller) CreateUser(ctx ligo.Context) error {
	var input CreateUserInput
	if err := ctx.Bind(&input); err != nil {
		return err
	}

	user, err := c.svc.CreateUser(input)
	if err != nil {
		return err
	}

	return ctx.JSON(201, user)
}

// UpdateUser handles PUT /users/:id
func (c *Controller) UpdateUser(ctx ligo.Context) error {
	id := ctx.Param("id")

	var input UpdateUserInput
	if err := ctx.Bind(&input); err != nil {
		return err
	}

	user, err := c.svc.UpdateUser(id, input)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.JSON(404, map[string]string{"error": "user not found"})
		}
		return err
	}

	return ctx.JSON(200, user)
}

// DeleteUser handles DELETE /users/:id
func (c *Controller) DeleteUser(ctx ligo.Context) error {
	id := ctx.Param("id")

	// Get current user from context (set by guard)
	currentUser, _ := ctx.Get(auth.ContextKeyUser).(*auth.User)

	c.log.Info("Deleting user",
		ligo.LoggerField{Key: "target_id", Value: id},
		ligo.LoggerField{Key: "actor_id", Value: currentUser.ID},
		ligo.LoggerField{Key: "actor_role", Value: currentUser.Role},
	)

	if err := c.svc.DeleteUser(id); err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.JSON(404, map[string]string{"error": "user not found"})
		}
		return err
	}

	return ctx.String(200, "")
}
