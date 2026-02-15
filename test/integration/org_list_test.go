package integration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/codemint/codemint-cli/internal/api"
)

func TestOrgList(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"organizations":[{"id":"o1","slug":"acme","name":"Acme","role":"owner"}]}`)),
		}, nil
	})

	c := api.NewClient(api.ClientOptions{BaseURL: "https://example.com", Timeout: 2 * time.Second, UserAgent: "test/1", Transport: transport})
	resp, err := c.OrgList(context.Background(), "token")
	if err != nil {
		t.Fatalf("OrgList error: %v", err)
	}
	if len(resp.Organizations) != 1 || resp.Organizations[0].Slug != "acme" {
		t.Fatalf("unexpected response: %+v", resp.Organizations)
	}
}

func TestOrgListArrayShape(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`[{"id":"o1","slug":"acme","name":"Acme","role":"owner"}]`)),
		}, nil
	})

	c := api.NewClient(api.ClientOptions{BaseURL: "https://example.com", Timeout: 2 * time.Second, UserAgent: "test/1", Transport: transport})
	resp, err := c.OrgList(context.Background(), "token")
	if err != nil {
		t.Fatalf("OrgList error: %v", err)
	}
	if len(resp.Organizations) != 1 || resp.Organizations[0].Slug != "acme" {
		t.Fatalf("unexpected response: %+v", resp.Organizations)
	}
}
