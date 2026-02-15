package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientOptions struct {
	BaseURL   string
	Timeout   time.Duration
	UserAgent string
	Debug     bool
	Transport http.RoundTripper
}

type Client struct {
	baseURL string
	http    *http.Client
	debug   bool
}

func NewClient(opts ClientOptions) *Client {
	tr := http.DefaultTransport
	if opts.Transport != nil {
		tr = opts.Transport
	}
	if opts.UserAgent != "" {
		tr = uaRoundTripper{next: tr, userAgent: opts.UserAgent}
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &Client{
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
		http:    &http.Client{Timeout: timeout, Transport: tr},
		debug:   opts.Debug,
	}
}

func (c *Client) AuthMe(ctx context.Context, token string) (*AuthMeResponse, error) {
	var out authMeEnvelope
	if err := c.do(ctx, http.MethodGet, "/api/auth/me", token, nil, &out); err != nil {
		return nil, err
	}
	return &out.User, nil
}

func (c *Client) ItemsSearch(ctx context.Context, token string, req ItemsSearchRequest) (*ItemsSearchResponse, error) {
	q := url.Values{}
	if req.Q != "" {
		q.Set("q", req.Q)
	}
	if req.Type != "" {
		q.Set("type", req.Type)
	}
	if req.Slug != "" {
		q.Set("slug", req.Slug)
	}
	if len(req.Tags) > 0 {
		q.Set("tags", strings.Join(req.Tags, ","))
	}
	if req.Latest {
		q.Set("latest", "true")
	}
	if req.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", req.Page))
	}
	if req.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", req.Limit))
	}
	path := "/api/items/search"
	if enc := q.Encode(); enc != "" {
		path += "?" + enc
	}
	var out ItemsSearchResponse
	if err := c.do(ctx, http.MethodGet, path, token, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CatalogSuggest(ctx context.Context, token string, req CatalogLookupRequest) ([]CatalogItem, error) {
	resp, err := c.ItemsSearch(ctx, token, ItemsSearchRequest{Q: req.Q, Type: req.Type, Tags: req.Tags, Latest: true, Page: 1, Limit: 50})
	if err != nil {
		return nil, err
	}
	out := make([]CatalogItem, 0, len(resp.Items))
	for _, it := range resp.Items {
		out = append(out, itemToCatalog(it))
	}
	return out, nil
}

func (c *Client) CatalogGetByRef(ctx context.Context, token string, itemType, slug string) (*CatalogItem, error) {
	ref := url.QueryEscape("@" + itemType + "/" + slug)
	var out CatalogItem
	if err := c.do(ctx, http.MethodGet, "/api/catalog/resolve?ref="+ref, token, nil, &out); err != nil {
		return nil, err
	}
	normalizeCatalogItem(&out)
	return &out, nil
}

func (c *Client) CatalogSync(ctx context.Context, token string, req CatalogSyncRequest) (*CatalogSyncResponse, error) {
	const maxBatch = 100
	apiItems := make([]*CatalogItem, 0, len(req.Items))
	for start := 0; start < len(req.Items); start += maxBatch {
		end := start + maxBatch
		if end > len(req.Items) {
			end = len(req.Items)
		}
		ids := make([]string, 0, end-start)
		for _, it := range req.Items[start:end] {
			ids = append(ids, it.CatalogID)
		}
		var apiOut catalogSyncAPIResponse
		if err := c.do(ctx, http.MethodPost, "/api/catalog/sync", token, map[string]any{"catalogIds": ids}, &apiOut); err != nil {
			return nil, err
		}
		apiItems = append(apiItems, apiOut.Items...)
	}

	results := make([]CatalogSyncResult, 0, len(req.Items))
	for i, local := range req.Items {
		var remote *CatalogItem
		if i < len(apiItems) {
			remote = apiItems[i]
		}
		if remote == nil {
			results = append(results, CatalogSyncResult{
				CatalogID:      local.CatalogID,
				CurrentVersion: local.Version,
				Removed:        true,
			})
			continue
		}
		normalizeCatalogItem(remote)
		results = append(results, CatalogSyncResult{
			CatalogID:      local.CatalogID,
			Slug:           remote.Slug,
			Type:           remote.Type,
			CurrentVersion: local.Version,
			LatestVersion:  remote.Version,
			Deprecated:     remote.Deprecated,
			Removed:        false,
			LatestItem:     *remote,
		})
	}
	return &CatalogSyncResponse{Results: results}, nil
}

func (c *Client) OrgList(ctx context.Context, token string) (*OrgListResponse, error) {
	var out OrgListResponse
	if err := c.do(ctx, http.MethodGet, "/api/org/my", token, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(ctx context.Context, method, path, token string, in any, out any) error {
	var payload []byte
	var err error
	if in != nil {
		payload, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}

	retries := 2
	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		var body io.Reader
		if len(payload) > 0 {
			body = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
		if err != nil {
			return err
		}
		if len(payload) > 0 {
			req.Header.Set("Content-Type", "application/json")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := c.http.Do(req)
		if err != nil {
			if transient(err) && attempt < retries {
				time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
				lastErr = err
				continue
			}
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			aerr := parseAPIError(resp.StatusCode, b)
			if resp.StatusCode >= 500 && attempt < retries {
				time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
				lastErr = aerr
				continue
			}
			return aerr
		}
		if out == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return lastErr
}

func parseAPIError(status int, b []byte) error {
	var env ErrorEnvelope
	if err := json.Unmarshal(b, &env); err == nil && env.Error.Message != "" {
		return &APIError{Status: status, Code: env.Error.Code, Message: env.Error.Message}
	}
	var flat map[string]any
	if err := json.Unmarshal(b, &flat); err == nil {
		if msg, ok := flat["error"].(string); ok && msg != "" {
			return &APIError{Status: status, Message: msg}
		}
		if msg, ok := flat["message"].(string); ok && msg != "" {
			return &APIError{Status: status, Message: msg}
		}
	}
	msg := strings.TrimSpace(string(b))
	if msg == "" {
		msg = http.StatusText(status)
	}
	return &APIError{Status: status, Message: msg}
}

func transient(err error) bool {
	if nerr, ok := err.(net.Error); ok {
		return nerr.Timeout() || nerr.Temporary()
	}
	return false
}

func itemToCatalog(it Item) CatalogItem {
	slug := it.Slug
	if slug == "" {
		slug = strMeta(it.Metadata, "slug")
	}
	catalogID := it.CatalogID
	if catalogID == "" {
		catalogID = strMeta(it.Metadata, "catalogId")
	}
	version := it.Version
	if version == "" {
		version = it.CatVer
	}
	if version == "" {
		version = strMeta(it.Metadata, "catalogVersion")
	}
	if version == "" {
		version = strMeta(it.Metadata, "version")
	}
	checksum := it.Checksum
	if checksum == "" {
		checksum = strMeta(it.Metadata, "checksum")
	}
	content := it.Content
	if content == "" {
		content = strMeta(it.Metadata, "content")
	}
	name := it.Name
	if name == "" {
		name = it.Title
	}
	if version == "" {
		version = "0.0.0"
	}
	if slug == "" {
		slug = it.ID
	}
	if catalogID == "" {
		catalogID = it.Type + ":" + slug
	}
	return CatalogItem{
		ID:         it.ID,
		Title:      it.Title,
		Name:       name,
		Type:       it.Type,
		Slug:       slug,
		CatalogID:  catalogID,
		Version:    version,
		CatVer:     version,
		Checksum:   checksum,
		Tags:       it.Tags,
		Deprecated: boolMeta(it.Metadata, "deprecated"),
		Changelog:  strMeta(it.Metadata, "changelog"),
		Content:    content,
		Metadata:   it.Metadata,
	}
}

func normalizeCatalogItem(it *CatalogItem) {
	if it.Name == "" {
		it.Name = it.Title
	}
	if it.Version == "" {
		it.Version = it.CatVer
	}
	if it.CatVer == "" {
		it.CatVer = it.Version
	}
	if it.CatalogID == "" && it.Type != "" && it.Slug != "" {
		it.CatalogID = it.Type + ":" + it.Slug
	}
	if it.Version == "" {
		it.Version = "0.0.0"
		it.CatVer = "0.0.0"
	}
}

func strMeta(meta map[string]any, key string) string {
	if meta == nil {
		return ""
	}
	v, ok := meta[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func boolMeta(meta map[string]any, key string) bool {
	if meta == nil {
		return false
	}
	v, ok := meta[key]
	if !ok || v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}
