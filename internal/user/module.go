package user

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/auth"
)

func Module() ligo.Module {
	return ligo.NewModule("user",
		ligo.Imports(auth.Module()),
		ligo.Providers(
			ligo.Export(ligo.Factory[*UserRepository](NewUserRepository)),
			ligo.Export(ligo.Factory[*UserService](NewUserService)),
		),
		ligo.Controllers(
			func(svc *UserService, log ligo.Logger, auth *auth.AuthService) ligo.Controller {
				return NewController(svc, log, auth)
			},
		),
	)
}
