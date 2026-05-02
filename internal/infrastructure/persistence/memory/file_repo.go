package memory

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

// FileRepository is an in-memory implementation of repository.FileRepository.
type FileRepository struct {
	mu    sync.RWMutex
	files map[string]*entity.File
	dir   string
}

// NewFileRepository creates a new in-memory file repository.
func NewFileRepository(dir string) repository.FileRepository {
	os.MkdirAll(dir, 0755)
	return &FileRepository{
		files: make(map[string]*entity.File),
		dir:   dir,
	}
}

func (r *FileRepository) Save(file io.Reader, filename string) (*entity.File, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	id := newUUID()
	path := filepath.Join(r.dir, id+"_"+filename)

	if err := os.WriteFile(path, content, 0644); err != nil {
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

	r.files[id] = fileEntity
	return fileEntity, nil
}

func (r *FileRepository) FindByID(id string) (*entity.File, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, found := r.files[id]
	return file, found
}

func (r *FileRepository) GetContent(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (r *FileRepository) FindAll() []*entity.File {
	r.mu.RLock()
	defer r.mu.RUnlock()

	files := make([]*entity.File, 0, len(r.files))
	for _, f := range r.files {
		files = append(files, f)
	}
	return files
}

func (r *FileRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, found := r.files[id]
	if !found {
		return fmt.Errorf("file not found")
	}

	if err := os.Remove(file.Path); err != nil {
		return err
	}

	delete(r.files, id)
	return nil
}

func detectContentType(filename string, content []byte) string {
	ct := http.DetectContentType(content)
	if ct != "application/octet-stream" {
		return ct
	}

	ext := filepath.Ext(filename)
	switch ext {
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
