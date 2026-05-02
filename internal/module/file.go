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

// File returns the file upload module.
func FileModule() ligo.Module {
	return ligo.NewModule("file",
		ligo.Providers(
			ligomemory.Provider[string, *entity.File](),
			ligo.Factory[repository.FileRepository](func(cfg *config.Config, store *ligomemory.Store[string, *entity.File]) repository.FileRepository {
				return memory.NewFileRepository(cfg.UploadDir, store)
			}),
			ligo.Factory[*usecase.FileUseCase](usecase.NewFileUseCase),
		),
		ligo.Controllers(controller.NewFileController),
	)
}
