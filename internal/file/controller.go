package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

const maxFileSize = 10 << 20 // 10MB

type Controller struct {
	log    ligo.Logger
	storage *FileStorage
}

func NewController(log ligo.Logger) ligo.Controller {
	uploadDir := filepath.Join(os.TempDir(), "ligo-uploads")
	os.MkdirAll(uploadDir, 0755)

	return &Controller{
		log:     log,
		storage: NewFileStorage(uploadDir),
	}
}

func (c *Controller) Routes(r ligo.Router) {
	cr := ligo.NewChainRouter(r.Group("/files"))

	cr.POST("/upload", c.Upload).
		Filter(common.GlobalExceptionFilter(c.log)).
		Handle()

	cr.GET("/:id", c.Download).
		Filter(common.GlobalExceptionFilter(c.log)).
		Handle()

	cr.GET("", c.ListFiles).
		Filter(common.GlobalExceptionFilter(c.log)).
		Handle()

	cr.DELETE("/:id", c.DeleteFile).
		Filter(common.GlobalExceptionFilter(c.log)).
		Handle()
}

func (c *Controller) Upload(ctx ligo.Context) error {
	if err := ctx.Request().ParseMultipartForm(maxFileSize); err != nil {
		return ctx.BadRequest("failed to parse form")
	}

	file, header, err := ctx.Request().FormFile("file")
	if err != nil {
		return ctx.BadRequest("file is required")
	}
	defer file.Close()

	savedFile, err := c.storage.Save(file, header.Filename)
	if err != nil {
		return ctx.InternalServerError("failed to save file")
	}

	c.log.Info("File uploaded",
		ligo.LoggerField{Key: "file_id", Value: savedFile.ID},
		ligo.LoggerField{Key: "filename", Value: header.Filename},
		ligo.LoggerField{Key: "size", Value: header.Size},
	)

	return ctx.Created(map[string]any{
		"id":       savedFile.ID,
		"filename": savedFile.Name,
		"size":     savedFile.Size,
		"url":      fmt.Sprintf("/files/%s", savedFile.ID),
	})
}

func (c *Controller) Download(ctx ligo.Context) error {
	id := ctx.Param("id")

	file, err := c.storage.Get(id)
	if err != nil {
		return err
	}
	defer file.Reader.Close()

	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	ctx.Response().Header().Set("Content-Type", file.ContentType)
	ctx.Response().Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))

	return ctx.Stream(file.Reader)
}

func (c *Controller) ListFiles(ctx ligo.Context) error {
	files := c.storage.List()

	fileInfos := make([]map[string]any, len(files))
	for i, f := range files {
		fileInfos[i] = map[string]any{
			"id":       f.ID,
			"filename": f.Name,
			"size":     f.Size,
			"url":      fmt.Sprintf("/files/%s", f.ID),
		}
	}

	return ctx.OK(map[string]any{
		"files": fileInfos,
		"count": len(fileInfos),
	})
}

func (c *Controller) DeleteFile(ctx ligo.Context) error {
	id := ctx.Param("id")

	if err := c.storage.Delete(id); err != nil {
		return err
	}

	c.log.Info("File deleted",
		ligo.LoggerField{Key: "file_id", Value: id},
	)

	return ctx.NoContent()
}

// File represents an uploaded file
type File struct {
	ID          string
	Name        string
	ContentType string
	Size        int64
	Path        string
	Reader      io.ReadCloser
}

// FileStorage handles file storage operations
type FileStorage struct {
	mu     sync.RWMutex
	files  map[string]*FileInfo
	nextID int
	dir    string
}

type FileInfo struct {
	ID          string
	Name        string
	ContentType string
	Size        int64
	Path        string
}

func NewFileStorage(dir string) *FileStorage {
	return &FileStorage{
		files:  make(map[string]*FileInfo),
		nextID: 1,
		dir:    dir,
	}
}

func (s *FileStorage) Save(file multipart.File, filename string) (*FileInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	id := strconv.Itoa(s.nextID)
	s.nextID++

	path := filepath.Join(s.dir, id+"_"+filename)

	if err := os.WriteFile(path, content, 0644); err != nil {
		return nil, err
	}

	contentType := detectContentType(filename, content)

	info := &FileInfo{
		ID:          id,
		Name:        filename,
		ContentType: contentType,
		Size:        int64(len(content)),
		Path:        path,
	}

	s.files[id] = info
	return info, nil
}

func (s *FileStorage) Get(id string) (*File, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, found := s.files[id]
	if !found {
		return nil, common.ErrNotFound
	}

	reader, err := os.Open(info.Path)
	if err != nil {
		return nil, err
	}

	return &File{
		ID:          info.ID,
		Name:        info.Name,
		ContentType: info.ContentType,
		Size:        info.Size,
		Path:        info.Path,
		Reader:      reader,
	}, nil
}

func (s *FileStorage) List() []*FileInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files := make([]*FileInfo, 0, len(s.files))
	for _, f := range s.files {
		files = append(files, f)
	}
	return files
}

func (s *FileStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, found := s.files[id]
	if !found {
		return common.ErrNotFound
	}

	if err := os.Remove(info.Path); err != nil {
		return err
	}

	delete(s.files, id)
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
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	default:
		return "application/octet-stream"
	}
}
