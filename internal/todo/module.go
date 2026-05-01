package todo

import (
	"github.com/linkeunid/ligo"
)

// Module returns the todo module with all providers and controllers
func Module() ligo.Module {
	return ligo.NewModule("todo",
		ligo.Providers(
			// Repository - singleton
			ligo.Factory[*Repository](NewRepository),
			// Service - singleton with auto-injected repo
			ligo.Factory[*Service](NewService),
		),
		// Controller constructor - dependencies auto-injected by framework
		ligo.Controllers(func(svc *Service) ligo.Controller {
			return NewController(svc)
		}),
	)
}
