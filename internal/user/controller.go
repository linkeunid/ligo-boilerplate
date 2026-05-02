package user

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/auth"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

type Controller struct {
	svc  *UserService
	log  ligo.Logger
	auth *auth.AuthService
}

func NewController(svc *UserService, log ligo.Logger, auth *auth.AuthService) ligo.Controller {
	return &Controller{svc: svc, log: log, auth: auth}
}

func (c *Controller) Routes(r ligo.Router) {
	cr := ligo.NewChainRouter(r.Group("/users"))

	cr.GET("", c.GetAllUsers).
		Filter(common.GlobalExceptionFilter(c.log)).
		Intercept(common.LoggingInterceptor(c.log)).
		Handle()

	cr.GET("/:id", c.GetUserByID).
		Filter(common.GlobalExceptionFilter(c.log)).
		Guard(auth.AuthGuard(c.auth)).
		Intercept(common.LoggingInterceptor(c.log)).
		Handle()

	cr.DELETE("/:id", c.DeleteUser).
		Filter(common.GlobalExceptionFilter(c.log)).
		Guard(auth.AuthGuard(c.auth), ligo.RolesGuard(auth.ContextKeyUser, "admin")).
		Intercept(common.LoggingInterceptor(c.log), auth.AuditInterceptor(c.log)).
		Handle()

	cr.POST("", c.CreateUser).
		Filter(common.GlobalExceptionFilter(c.log)).
		Guard(auth.AuthGuard(c.auth)).
		Intercept(common.LoggingInterceptor(c.log)).
		Handle()

	cr.PUT("/:id", c.UpdateUser).
		Filter(common.GlobalExceptionFilter(c.log)).
		Guard(auth.AuthGuard(c.auth)).
		Intercept(common.LoggingInterceptor(c.log)).
		Handle()
}

func (c *Controller) GetAllUsers(ctx ligo.Context) error {
	users := c.svc.GetAllUsers()
	return ctx.OK(map[string]any{
		"users": users,
		"count": len(users),
	})
}

func (c *Controller) GetUserByID(ctx ligo.Context) error {
	id := ctx.Param("id")
	user, err := c.svc.GetUserByID(id)
	if err != nil {
		return err
	}
	return ctx.OK(user)
}

func (c *Controller) CreateUser(ctx ligo.Context) error {
	var input CreateUserInput
	if err := ctx.Bind(&input); err != nil {
		return common.ErrValidation
	}
	user, err := c.svc.CreateUser(input)
	if err != nil {
		return err
	}
	return ctx.Created(user)
}

func (c *Controller) UpdateUser(ctx ligo.Context) error {
	id := ctx.Param("id")
	var input UpdateUserInput
	if err := ctx.Bind(&input); err != nil {
		return common.ErrValidation
	}
	user, err := c.svc.UpdateUser(id, input)
	if err != nil {
		return err
	}
	return ctx.OK(user)
}

func (c *Controller) DeleteUser(ctx ligo.Context) error {
	id := ctx.Param("id")

	if err := c.svc.DeleteUser(id); err != nil {
		return err
	}
	return ctx.NoContent()
}
