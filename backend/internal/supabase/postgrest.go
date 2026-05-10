package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ErrNotFound is returned when a single-row query returns no rows.
var ErrNotFound = errors.New("not found")

// Filters is a thin alias around url.Values for PostgREST query strings.
//
// Examples:
//   filters.Set("id", "eq.<uuid>")
//   filters.Set("select", "id,title,owner_id")
//   filters.Set("order", "created_at.desc")
//   filters.Set("limit", "20")
type Filters = url.Values

func NewFilters() Filters { return url.Values{} }

func (c *Client) restURL(table string) string {
	return c.baseURL + "/rest/v1/" + table
}

// Select runs a GET against /rest/v1/<table> and decodes the JSON array into dst.
func (c *Client) Select(ctx context.Context, table string, f Filters, dst any) error {
	u := c.restURL(table)
	if len(f) > 0 {
		u += "?" + f.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	c.applyAuth(req, false, "")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decode(resp, dst)
}

// SelectOne fetches exactly one row. Returns ErrNotFound if the result is empty.
// dst must be a pointer to the row struct.
func (c *Client) SelectOne(ctx context.Context, table string, f Filters, dst any) error {
	if f == nil {
		f = NewFilters()
	}
	f.Set("limit", "1")
	u := c.restURL(table) + "?" + f.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.pgrst.object+json")
	c.applyAuth(req, false, "")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotAcceptable || resp.StatusCode == 404 ||
		(resp.StatusCode == 200 && len(body) == 0) {
		return ErrNotFound
	}
	if resp.StatusCode >= 400 {
		// PostgREST returns PGRST116 when expecting a single row and 0 are found.
		var apiErr APIError
		_ = json.Unmarshal(body, &apiErr)
		apiErr.Status = resp.StatusCode
		apiErr.Raw = string(body)
		if strings.Contains(apiErr.Raw, "PGRST116") {
			return ErrNotFound
		}
		return &apiErr
	}
	if dst == nil {
		return nil
	}
	return json.Unmarshal(body, dst)
}

// Insert posts one row (or many) into <table>. If `returning` is non-nil, the
// inserted row(s) are decoded into it (Prefer: return=representation).
func (c *Client) Insert(ctx context.Context, table string, body any, returning any) error {
	return c.write(ctx, http.MethodPost, c.restURL(table), body, "", returning)
}

// Upsert posts rows with on-conflict resolution. `onConflict` is a comma-separated
// column list matching a unique index (e.g. "video_id,quality").
func (c *Client) Upsert(ctx context.Context, table string, body any, onConflict string, returning any) error {
	u := c.restURL(table)
	if onConflict != "" {
		u += "?on_conflict=" + onConflict
	}
	return c.write(ctx, http.MethodPost, u, body, "resolution=merge-duplicates", returning)
}

// Update applies a PATCH with the given filters. `returning` is optional.
func (c *Client) Update(ctx context.Context, table string, f Filters, body any, returning any) error {
	if len(f) == 0 {
		return errors.New("update without filters is forbidden")
	}
	u := c.restURL(table) + "?" + f.Encode()
	return c.write(ctx, http.MethodPatch, u, body, "", returning)
}

// Delete removes rows matching the filters and returns the number of deleted rows.
func (c *Client) Delete(ctx context.Context, table string, f Filters) (int, error) {
	if len(f) == 0 {
		return 0, errors.New("delete without filters is forbidden")
	}
	u := c.restURL(table) + "?" + f.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Prefer", "return=representation")
	c.applyAuth(req, false, "")
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	var rows []json.RawMessage
	if err := decode(resp, &rows); err != nil {
		return 0, err
	}
	return len(rows), nil
}

// Count returns the total row count for the given filters using Prefer: count=exact.
func (c *Client) Count(ctx context.Context, table string, f Filters) (int64, error) {
	u := c.restURL(table) + "?select=*"
	if len(f) > 0 {
		u += "&" + f.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Prefer", "count=exact")
	req.Header.Set("Range-Unit", "items")
	req.Header.Set("Range", "0-0")
	c.applyAuth(req, false, "")
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		return 0, &APIError{Status: resp.StatusCode, Raw: string(body)}
	}
	cr := resp.Header.Get("Content-Range")
	if cr == "" {
		return 0, nil
	}
	idx := strings.Index(cr, "/")
	if idx < 0 {
		return 0, nil
	}
	tail := cr[idx+1:]
	if tail == "*" {
		return 0, nil
	}
	n, err := strconv.ParseInt(tail, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Exists is true if at least one row matches.
func (c *Client) Exists(ctx context.Context, table string, f Filters) (bool, error) {
	n, err := c.Count(ctx, table, f)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// RPC calls /rest/v1/rpc/<fn> with the given JSON body.
func (c *Client) RPC(ctx context.Context, fn string, body any, returning any) error {
	return c.write(ctx, http.MethodPost, c.baseURL+"/rest/v1/rpc/"+fn, body, "", returning)
}

// PingREST hits the PostgREST root for a quick liveness check.
func (c *Client) PingREST(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/rest/v1/", nil)
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
	return nil
}

// write is the shared body for POST/PATCH PostgREST calls.
func (c *Client) write(ctx context.Context, method, u string, body any, extraPrefer string, returning any) error {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	prefer := ""
	if returning != nil {
		prefer = "return=representation"
	} else {
		prefer = "return=minimal"
	}
	if extraPrefer != "" {
		prefer = extraPrefer + "," + prefer
	}
	req.Header.Set("Prefer", prefer)
	c.applyAuth(req, false, "")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if returning == nil {
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			apiErr := &APIError{Status: resp.StatusCode, Raw: string(body)}
			_ = json.Unmarshal(body, apiErr)
			return apiErr
		}
		return nil
	}
	return decode(resp, returning)
}
