package module

import (
	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/persistence/memory"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

// User returns the user module.
func UserModule() ligo.Module {
	return ligo.NewModule("user",
		ligo.Providers(
			ligomemory.Provider[string, *entity.User](),
			ligo.Factory[repository.UserRepository](memory.NewUserRepository),
			ligo.Factory[*usecase.UserUseCase](usecase.NewUserUseCase),
		),
		ligo.Controllers(controller.NewUserController),
	)
}
