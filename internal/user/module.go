package user

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/auth"
)

// Module returns the user module with all providers and controllers.
// This demonstrates a complete module with all Ligo features.
func Module() ligo.Module {
	return ligo.NewModule("user",
		// Import auth module to access AuthService
		ligo.Imports(auth.Module()),
		// Providers
		ligo.Providers(
			// User Repository - singleton
			ligo.Export(ligo.Factory[*UserRepository](NewUserRepository)),
			// User Service - singleton with auto-injected repo and logger
			ligo.Export(ligo.Factory[*UserService](NewUserService)),
		),
		// Controllers
		ligo.Controllers(
			// User Controller - auto-injected with service, logger, and auth service
			func(svc *UserService, log ligo.Logger, auth *auth.AuthService) ligo.Controller {
				return NewController(svc, log, auth)
			},
		),
	)
}
