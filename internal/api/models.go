package api

import "time"

type AuthMeResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type authMeEnvelope struct {
	User AuthMeResponse `json:"user"`
}

type Item struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Title     string                 `json:"title,omitempty"`
	Type      string                 `json:"type"`
	Slug      string                 `json:"slug,omitempty"`
	CatalogID string                 `json:"catalogId,omitempty"`
	Version   string                 `json:"version,omitempty"`
	CatVer    string                 `json:"catalogVersion,omitempty"`
	Checksum  string                 `json:"checksum,omitempty"`
	Tags      []string               `json:"tags"`
	Score     int                    `json:"score"`
	Metadata  map[string]any         `json:"metadata,omitempty"`
	Content   string                 `json:"content,omitempty"`
	UpdatedAt string                 `json:"updatedAt,omitempty"`
	CreatedAt string                 `json:"createdAt,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

type ItemsSearchRequest struct {
	Q      string   `json:"q"`
	Type   string   `json:"type,omitempty"`
	Slug   string   `json:"slug,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Latest bool     `json:"latest,omitempty"`
	Page   int      `json:"page,omitempty"`
	Limit  int      `json:"limit,omitempty"`
}

type ItemsSearchResponse struct {
	Items []Item `json:"items"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}

type CatalogItem struct {
	ID         string         `json:"id"`
	Title      string         `json:"title,omitempty"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Slug       string         `json:"slug"`
	CatalogID  string         `json:"catalogId"`
	Version    string         `json:"version,omitempty"`
	CatVer     string         `json:"catalogVersion,omitempty"`
	Checksum   string         `json:"checksum"`
	Tags       []string       `json:"tags"`
	Deprecated bool           `json:"deprecated,omitempty"`
	Changelog  string         `json:"changelog,omitempty"`
	Content    string         `json:"content,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type CatalogLookupRequest struct {
	Type string   `json:"type,omitempty"`
	Tags []string `json:"tags,omitempty"`
	Q    string   `json:"q,omitempty"`
}

type CatalogSyncRequest struct {
	Items []CatalogSyncItem `json:"items"`
}

type CatalogSyncItem struct {
	CatalogID string `json:"catalogId"`
	Version   string `json:"version"`
	Checksum  string `json:"checksum"`
}

type CatalogSyncResult struct {
	CatalogID      string      `json:"catalogId"`
	Slug           string      `json:"slug"`
	Type           string      `json:"type"`
	CurrentVersion string      `json:"currentVersion"`
	LatestVersion  string      `json:"latestVersion"`
	Deprecated     bool        `json:"deprecated"`
	Removed        bool        `json:"removed"`
	LatestItem     CatalogItem `json:"latestItem"`
}

type CatalogSyncResponse struct {
	Results []CatalogSyncResult `json:"results"`
}

type catalogSyncAPIResponse struct {
	Items []*CatalogItem `json:"items"`
}

type Organization struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type OrgListResponse struct {
	Organizations []Organization `json:"organizations"`
}

type CLIAuthCallbackPayload struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

type InstallRecord struct {
	CatalogID   string    `json:"catalogId"`
	Type        string    `json:"type"`
	Slug        string    `json:"slug"`
	Version     string    `json:"version"`
	Checksum    string    `json:"checksum"`
	InstalledAt time.Time `json:"installedAt"`
}
