package handlers

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"library_back/internal/pkg/cloudinary"

	"github.com/gin-gonic/gin"
)

type stubUploader struct {
	lastFolder string
	result     *cloudinary.UploadResult
	err        error
}

func (s *stubUploader) UploadImage(_ context.Context, folder string, _ string, _ io.Reader, _ int64) (*cloudinary.UploadResult, error) {
	s.lastFolder = folder
	if s.err != nil {
		return nil, s.err
	}
	return s.result, nil
}

func TestUpload_Create201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	up := &stubUploader{result: &cloudinary.UploadResult{
		URL: "https://res.cloudinary.com/demo/image/upload/v1/lumina/books/x.jpg", PublicID: "lumina/books/x", Folder: "lumina/books",
	}}
	h := NewUploadHandler(up)
	r := gin.New()
	r.POST("/api/uploads", h.Create)

	body, ctype := multipartImage(t, "cover.jpg", tinyJPEG(), "books")
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", body)
	req.Header.Set("Content-Type", ctype)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if up.lastFolder != "books" {
		t.Fatalf("folder=%s", up.lastFolder)
	}
}

func TestUpload_RejectsNonImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler(&stubUploader{result: &cloudinary.UploadResult{URL: "x"}})
	r := gin.New()
	r.POST("/api/uploads", h.Create)

	body, ctype := multipartImage(t, "note.txt", []byte("hello world not an image"), "books")
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", body)
	req.Header.Set("Content-Type", ctype)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestUpload_NotConfigured503(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler(nil)
	r := gin.New()
	r.POST("/api/uploads", h.Create)

	body, ctype := multipartImage(t, "cover.jpg", tinyJPEG(), "books")
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", body)
	req.Header.Set("Content-Type", ctype)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", w.Code)
	}
}

func multipartImage(t *testing.T, filename string, data []byte, folder string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("folder", folder)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf, w.FormDataContentType()
}

// Minimal valid JPEG (1x1) for content-type sniffing.
func tinyJPEG() []byte {
	return []byte{
		0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43,
		0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08, 0x07, 0x07, 0x07, 0x09,
		0x09, 0x08, 0x0a, 0x0c, 0x14, 0x0d, 0x0c, 0x0b, 0x0b, 0x0c, 0x19, 0x12,
		0x13, 0x0f, 0x14, 0x1d, 0x1a, 0x1f, 0x1e, 0x1d, 0x1a, 0x1c, 0x1c, 0x20,
		0x24, 0x2e, 0x27, 0x20, 0x22, 0x2c, 0x23, 0x1c, 0x1c, 0x28, 0x37, 0x29,
		0x2c, 0x30, 0x31, 0x34, 0x34, 0x34, 0x1f, 0x27, 0x39, 0x3d, 0x38, 0x32,
		0x3c, 0x2e, 0x33, 0x34, 0x32, 0xff, 0xc0, 0x00, 0x0b, 0x08, 0x00, 0x01,
		0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xff, 0xc4, 0x00, 0x14, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x03, 0xff, 0xc4, 0x00, 0x14, 0x10, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xda, 0x00, 0x08, 0x01, 0x01, 0x00, 0x00, 0x3f, 0x00,
		0x37, 0xff, 0xd9,
	}
}
