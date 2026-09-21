package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// MaxImageBytes is the maximum accepted upload size (5 MiB).
	MaxImageBytes = 5 << 20
	defaultAPIBase = "https://api.cloudinary.com/v1_1"
)

// AllowedFolders restricts where staff can store assets.
var AllowedFolders = map[string]struct{}{
	"books":     {},
	"students":  {},
	"badges":    {},
	"users":     {},
}

// Uploader uploads image bytes to an external media provider.
type Uploader interface {
	UploadImage(ctx context.Context, folder string, filename string, r io.Reader, size int64) (*UploadResult, error)
}

// UploadResult is the public URL payload returned to API clients.
type UploadResult struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
	Folder   string `json:"folder"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Format   string `json:"format,omitempty"`
	Bytes    int    `json:"bytes,omitempty"`
}

// Client talks to the Cloudinary Upload API with signed requests.
type Client struct {
	cloudName string
	apiKey    string
	apiSecret string
	folderRoot string
	httpClient *http.Client
	apiBase    string
	now        func() time.Time
}

// ConfigFromEnv builds a Client from CLOUDINARY_* environment variables.
// Returns nil, error if required vars are missing.
func ConfigFromEnv() (*Client, error) {
	cloud := strings.TrimSpace(os.Getenv("CLOUDINARY_CLOUD_NAME"))
	key := strings.TrimSpace(os.Getenv("CLOUDINARY_API_KEY"))
	secret := strings.TrimSpace(os.Getenv("CLOUDINARY_API_SECRET"))
	if cloud == "" || key == "" || secret == "" {
		return nil, fmt.Errorf("cloudinary: CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY and CLOUDINARY_API_SECRET are required")
	}
	root := strings.Trim(strings.TrimSpace(os.Getenv("CLOUDINARY_FOLDER_ROOT")), "/")
	if root == "" {
		root = "lumina"
	}
	return &Client{
		cloudName:  cloud,
		apiKey:     key,
		apiSecret:  secret,
		folderRoot: root,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		apiBase:    defaultAPIBase,
		now:        time.Now,
	}, nil
}

// UploadImage streams an image to Cloudinary under folderRoot/<folder>.
func (c *Client) UploadImage(ctx context.Context, folder string, filename string, r io.Reader, size int64) (*UploadResult, error) {
	folder = strings.Trim(strings.ToLower(strings.TrimSpace(folder)), "/")
	if _, ok := AllowedFolders[folder]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidFolder, folder)
	}
	if size <= 0 {
		return nil, ErrEmptyFile
	}
	if size > MaxImageBytes {
		return nil, ErrTooLarge
	}

	targetFolder := path.Join(c.folderRoot, folder)
	ts := strconv.FormatInt(c.now().Unix(), 10)
	params := map[string]string{
		"folder":    targetFolder,
		"timestamp": ts,
	}
	signature := signParams(params, c.apiSecret)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("api_key", c.apiKey)
	_ = w.WriteField("timestamp", ts)
	_ = w.WriteField("folder", targetFolder)
	_ = w.WriteField("signature", signature)

	part, err := w.CreateFormFile("file", sanitizeFilename(filename))
	if err != nil {
		return nil, fmt.Errorf("cloudinary form file: %w", err)
	}
	limited := io.LimitReader(r, MaxImageBytes+1)
	written, err := io.Copy(part, limited)
	if err != nil {
		return nil, fmt.Errorf("cloudinary copy: %w", err)
	}
	if written > MaxImageBytes {
		return nil, ErrTooLarge
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("cloudinary form close: %w", err)
	}

	url := fmt.Sprintf("%s/%s/image/upload", c.apiBase, c.cloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudinary request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d: %s", ErrProvider, res.StatusCode, truncate(string(raw), 200))
	}

	var parsed cloudinaryResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("cloudinary decode: %w", err)
	}
	if parsed.SecureURL == "" {
		return nil, fmt.Errorf("%w: empty secure_url", ErrProvider)
	}
	return &UploadResult{
		URL:      parsed.SecureURL,
		PublicID: parsed.PublicID,
		Folder:   targetFolder,
		Width:    parsed.Width,
		Height:   parsed.Height,
		Format:   parsed.Format,
		Bytes:    parsed.Bytes,
	}, nil
}

type cloudinaryResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
	Bytes     int    `json:"bytes"`
}

func signParams(params map[string]string, apiSecret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	payload := strings.Join(parts, "&") + apiSecret
	sum := sha1.Sum([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func sanitizeFilename(name string) string {
	name = path.Base(strings.TrimSpace(name))
	if name == "" || name == "." || name == "/" {
		return "upload.bin"
	}
	return name
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
