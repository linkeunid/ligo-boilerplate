package module

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/config"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/persistence/memory"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

// File returns the file upload module.
func File() ligo.Module {
	return ligo.NewModule("file",
		ligo.Providers(
			ligo.Factory[repository.FileRepository](func(cfg *config.Config) repository.FileRepository {
				return memory.NewFileRepository(cfg.UploadDir)
			}),
			ligo.Factory[*usecase.FileUseCase](usecase.NewFileUseCase),
		),
		ligo.Controllers(controller.NewFileController),
	)
}
