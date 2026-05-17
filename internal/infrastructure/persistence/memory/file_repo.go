package memory

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ligomemory "github.com/linkeunid/ligo-memory"

	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

// FileRepository is an in-memory implementation of repository.FileRepository.
// Metadata is stored in a ligo-memory Store; file content is written to disk.
type FileRepository struct {
	store *ligomemory.Store[int, *entity.File]
	dir   string
}

// NewFileRepository creates a new in-memory file repository.
func NewFileRepository(dir string, store *ligomemory.Store[int, *entity.File]) repository.FileRepository {
	return &FileRepository{store: store, dir: dir}
}

// OnModuleInit initializes the file repository by creating the upload directory.
func (r *FileRepository) OnModuleInit() error {
	if err := os.MkdirAll(r.dir, 0o755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}
	return nil
}

// OnModuleDestroy cleans up resources when the module is destroyed.
// This is called in reverse order during shutdown.
func (r *FileRepository) OnModuleDestroy() error {
	// Optional: Cleanup uploaded files on shutdown
	// Uncomment if you want to clean up on shutdown:
	// return os.RemoveAll(r.dir)
	return nil
}

func (r *FileRepository) Save(file io.Reader, filename string) (*entity.File, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	id := nextID()
	path := filepath.Join(r.dir, fmt.Sprintf("%d_%s", id, filename))

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return nil, err
	}

	fileEntity := &entity.File{
		ID:          id,
		Name:        filename,
		ContentType: detectContentType(filename, content),
		Size:        int64(len(content)),
		Path:        path,
		CreatedAt:   time.Now(),
	}

	r.store.Set(id, fileEntity)
	return fileEntity, nil
}

func (r *FileRepository) FindByID(id int) (*entity.File, bool) {
	return r.store.Get(id)
}

func (r *FileRepository) GetContent(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (r *FileRepository) FindAll() []*entity.File {
	return r.store.All()
}

func (r *FileRepository) Delete(id int) error {
	file, found := r.store.Get(id)
	if !found {
		return fmt.Errorf("file not found")
	}
	if err := os.Remove(file.Path); err != nil {
		return err
	}
	r.store.Delete(id)
	return nil
}

func detectContentType(filename string, content []byte) string {
	ct := http.DetectContentType(content)
	if ct != "application/octet-stream" {
		return ct
	}
	switch filepath.Ext(filename) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
