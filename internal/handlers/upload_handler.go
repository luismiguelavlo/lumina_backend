package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"library_back/internal/models"
	"library_back/internal/pkg/cloudinary"

	"github.com/gin-gonic/gin"
)

const (
	msgUploadInvalid  = "Archivo inválido"
	msgUploadTooLarge = "El archivo supera el tamaño máximo (5 MB)"
	msgUploadFolder   = "Carpeta de subida inválida"
	msgUploadMime     = "Tipo de archivo no permitido (usa JPEG, PNG, WebP o GIF)"
	msgUploadNotReady = "Subida de imágenes no configurada en el servidor"
	msgUploadFailed   = "No se pudo subir la imagen"
)

var allowedImageMIME = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

// UploadHandler exposes POST /api/uploads (staff only).
type UploadHandler struct {
	uploader cloudinary.Uploader
}

// NewUploadHandler constructs UploadHandler. uploader may be nil → 503.
func NewUploadHandler(uploader cloudinary.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

// UploadResponse is returned on successful upload.
type UploadResponse struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
	Folder   string `json:"folder"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Format   string `json:"format,omitempty"`
	Bytes    int    `json:"bytes,omitempty"`
}

// Create POST /api/uploads
// multipart form: file (required), folder (optional: books|students|badges|users; default books)
func (h *UploadHandler) Create(c *gin.Context) {
	if h.uploader == nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{Message: msgUploadNotReady})
		return
	}

	folder := strings.TrimSpace(strings.ToLower(c.DefaultPostForm("folder", "books")))
	if _, ok := cloudinary.AllowedFolders[folder]; !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgUploadFolder,
			Errors:  map[string]string{"folder": "Debe ser books, students, badges o users"},
		})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgUploadInvalid,
			Errors:  map[string]string{"file": "Campo file requerido (multipart)"},
		})
		return
	}
	if fileHeader.Size <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgUploadInvalid})
		return
	}
	if fileHeader.Size > cloudinary.MaxImageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{Message: msgUploadTooLarge})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgUploadInvalid})
		return
	}
	defer func() { _ = src.Close() }()

	payload, err := io.ReadAll(io.LimitReader(src, cloudinary.MaxImageBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgUploadInvalid})
		return
	}
	if int64(len(payload)) > cloudinary.MaxImageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{Message: msgUploadTooLarge})
		return
	}

	mime := http.DetectContentType(payload)
	if _, ok := allowedImageMIME[mime]; !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgUploadMime,
			Errors:  map[string]string{"file": mime},
		})
		return
	}

	result, err := h.uploader.UploadImage(
		c.Request.Context(),
		folder,
		fileHeader.Filename,
		bytes.NewReader(payload),
		int64(len(payload)),
	)
	if err != nil {
		writeUploadError(c, err)
		return
	}
	c.JSON(http.StatusCreated, UploadResponse{
		URL:      result.URL,
		PublicID: result.PublicID,
		Folder:   result.Folder,
		Width:    result.Width,
		Height:   result.Height,
		Format:   result.Format,
		Bytes:    result.Bytes,
	})
}

func writeUploadError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, cloudinary.ErrInvalidFolder):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgUploadFolder})
	case errors.Is(err, cloudinary.ErrTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{Message: msgUploadTooLarge})
	case errors.Is(err, cloudinary.ErrEmptyFile):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgUploadInvalid})
	case errors.Is(err, cloudinary.ErrNotConfigured):
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{Message: msgUploadNotReady})
	default:
		c.JSON(http.StatusBadGateway, models.ErrorResponse{Message: msgUploadFailed})
	}
}
