package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SignedUploadURL struct {
	URL       string `json:"url"`
	Token     string `json:"token"`
	Path      string `json:"-"`
	ExpiresAt int64  `json:"-"`
}

// CreateSignedUploadURL returns a one-shot signed upload URL for a path.
// The frontend uploads directly to this URL with PUT.
func (c *Client) CreateSignedUploadURL(bucket, path string) (*SignedUploadURL, error) {
	endpoint := "/storage/v1/object/upload/sign/" + bucket + "/" + escapePath(path)
	resp, err := c.do(http.MethodPost, endpoint, map[string]any{}, false, "")
	if err != nil {
		return nil, err
	}
	var out SignedUploadURL
	if err := decode(resp, &out); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(out.URL, "http") {
		if strings.HasPrefix(out.URL, "/storage/v1/") {
			out.URL = c.baseURL + out.URL
		} else {
			out.URL = c.baseURL + "/storage/v1" + out.URL
		}
	}
	out.Path = path
	out.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
	return &out, nil
}

// PublicURL builds the public CDN URL for an object in a public bucket.
func (c *Client) PublicURL(bucket, path string) string {
	return c.baseURL + "/storage/v1/object/public/" + bucket + "/" + path
}

// Download writes a stored object's bytes to dstPath. Streams via io.Copy
// so very large videos do not buffer in memory.
func (c *Client) Download(ctx context.Context, bucket, path, dstPath string) error {
	endpoint := "/storage/v1/object/" + bucket + "/" + escapePath(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+endpoint, nil)
	if err != nil {
		return err
	}
	c.applyAuth(req, false, "")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return &APIError{Status: resp.StatusCode, Raw: string(body)}
	}
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

// UploadFile uploads a local file to bucket/path. Existing objects are overwritten.
func (c *Client) UploadFile(ctx context.Context, bucket, path, srcPath string) (int64, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	contentType := guessContentType(srcPath)

	endpoint := "/storage/v1/object/" + bucket + "/" + escapePath(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, f)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")
	req.Header.Set("Cache-Control", "public, max-age=3600")
	req.ContentLength = stat.Size()
	c.applyAuth(req, false, "")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return 0, &APIError{Status: resp.StatusCode, Raw: string(body)}
	}
	return stat.Size(), nil
}

// Remove deletes an object.
func (c *Client) Remove(ctx context.Context, bucket, path string) error {
	endpoint := "/storage/v1/object/" + bucket + "/" + escapePath(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+endpoint, nil)
	if err != nil {
		return err
	}
	c.applyAuth(req, false, "")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}

// StorageObject is one entry returned by List.
type StorageObject struct {
	Name           string         `json:"name"`
	ID             string         `json:"id,omitempty"`
	UpdatedAt      string         `json:"updated_at,omitempty"`
	CreatedAt      string         `json:"created_at,omitempty"`
	LastAccessedAt string         `json:"last_accessed_at,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// Size extracts the size in bytes from the metadata blob.
func (o *StorageObject) Size() int64 {
	if o.Metadata == nil {
		return 0
	}
	switch v := o.Metadata["size"].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	}
	return 0
}

type BucketUsage struct {
	Bucket    string `json:"bucket"`
	Files     int64  `json:"files"`
	SizeBytes int64  `json:"size_bytes"`
}

// BucketUsage walks the bucket root and aggregates file count + total size.
// Designed for the admin dashboard; for very large buckets, replace with a
// background job that caches results.
func (c *Client) BucketUsage(ctx context.Context, bucket string) (*BucketUsage, error) {
	const limit = 100
	usage := &BucketUsage{Bucket: bucket}

	for offset := 0; ; offset += limit {
		body := map[string]any{
			"limit":  limit,
			"offset": offset,
			"prefix": "",
		}
		endpoint := "/storage/v1/object/list/" + bucket
		req, err := newJSONReq(ctx, http.MethodPost, c.baseURL+endpoint, body)
		if err != nil {
			return nil, err
		}
		c.applyAuth(req, false, "")
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		var page []StorageObject
		if err := decode(resp, &page); err != nil {
			return nil, err
		}
		if len(page) == 0 {
			break
		}
		for _, o := range page {
			usage.Files++
			usage.SizeBytes += o.Size()
		}
		if len(page) < limit {
			break
		}
	}
	return usage, nil
}

// --- helpers ---

func escapePath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func guessContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/mp2t"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	}
	if ct := mime.TypeByExtension(filepath.Ext(name)); ct != "" {
		return ct
	}
	return "application/octet-stream"
}

func newJSONReq(ctx context.Context, method, urlStr string, body any) (*http.Request, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
