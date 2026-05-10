// Package supabase provides a thin HTTP client for the Supabase GoTrue (Auth)
// and Storage REST APIs. API keys are optional so self-hosted deployments that
// do not require Kong apikey headers can run without them.
package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	serviceKey string
	anonKey    string
	http       *http.Client
}

func New(baseURL, serviceKey, anonKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		serviceKey: serviceKey,
		anonKey:    anonKey,
		http: &http.Client{
			// Long timeout so worker downloads of large videos don't drop. Worker also
			// wraps calls in a context with its own deadline.
			Timeout: 0,
		},
	}
}

func (c *Client) HasServiceKey() bool {
	return c.serviceKey != ""
}

type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"error_code,omitempty"`
	Message string `json:"msg,omitempty"`
	Raw     string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("supabase: %d %s", e.Status, e.Message)
	}
	return fmt.Sprintf("supabase: %d %s", e.Status, e.Raw)
}

// do performs a JSON request. If a service-role key is configured it is used
// for backend calls; otherwise the anon key is used as the single project API
// key. If no key is configured, no apikey/bearer key headers are sent.
func (c *Client) do(method, path string, body any, useAnon bool, userToken string) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.baseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req, useAnon, userToken)
	return c.http.Do(req)
}

// doRaw performs a request with an arbitrary body (used for storage uploads).
// timeout=0 means no per-request timeout (caller must use ctx).
func (c *Client) doRaw(method, path string, body io.Reader, contentType string, headers map[string]string, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	c.applyAuth(req, false, "")
	if timeout > 0 {
		client := &http.Client{Timeout: timeout}
		return client.Do(req)
	}
	return c.http.Do(req)
}

func (c *Client) applyAuth(req *http.Request, useAnon bool, userToken string) {
	apiKey := c.serviceKey
	if apiKey == "" {
		apiKey = c.anonKey
	}
	if useAnon && c.anonKey != "" {
		apiKey = c.anonKey
	}
	if userToken != "" {
		req.Header.Set("Authorization", "Bearer "+userToken)
	}
	if apiKey == "" {
		return
	}
	req.Header.Set("apikey", apiKey)
	if userToken == "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
}

func decode(resp *http.Response, dst any) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		apiErr := &APIError{Status: resp.StatusCode, Raw: string(body)}
		_ = json.Unmarshal(body, apiErr)
		return apiErr
	}
	if dst == nil || len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, dst)
}
