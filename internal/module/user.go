package module

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/persistence/memory"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

// User returns the user module.
func User() ligo.Module {
	return ligo.NewModule("user",
		ligo.Providers(
			ligo.Factory[repository.UserRepository](memory.NewUserRepository),
			ligo.Factory[*usecase.UserUseCase](usecase.NewUserUseCase),
		),
		ligo.Controllers(controller.NewUserController),
	)
}
