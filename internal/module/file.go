package module

import (
	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"

	"github.com/linkeunid/ligo-boilerplate/internal/config"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/persistence/memory"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

// FileModule returns the file upload module with lifecycle hooks.
func FileModule() ligo.Module {
	return ligo.NewModule(
		"file",
		ligo.Providers(
			ligomemory.Provider[int, *entity.File](),
			ligo.Factory[repository.FileRepository](func(cfg *config.Config, store *ligomemory.Store[int, *entity.File], log ligo.Logger) repository.FileRepository {
				return memory.NewFileRepository(cfg.UploadDir, store)
			}),
			ligo.Factory[*usecase.FileUseCase](usecase.NewFileUseCase),
		),
		// Use HookedController for compile-time safe hook registration.
		// The FileController.Register method will be called automatically
		// to register its lifecycle hooks.
		ligo.Controllers(ligo.HookedController(controller.NewFileController)),
		// Module-level init hook: runs when module initializes
		ligo.OnModuleInit(func() error {
			// Module-level hooks don't support logger injection
			// Use provider-level hooks for logging, or get logger via DI
			return nil
		}),
		// Module-level destroy hook: runs when module destroys (reverse order)
		ligo.OnModuleDestroy(func() error {
			// Module-level hooks don't support logger injection
			// Use provider-level hooks for logging, or get logger via DI
			return nil
		}),
	)
}
