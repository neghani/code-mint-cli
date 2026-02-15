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

func TestItemsSearch(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Get("q") != "alpha" {
			t.Fatalf("missing query")
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"items":[{"id":"i1","name":"Alpha","type":"pkg","tags":["x"],"score":10}],"page":1,"limit":20,"total":1}`)),
		}, nil
	})

	c := api.NewClient(api.ClientOptions{BaseURL: "https://example.com", Timeout: 2 * time.Second, UserAgent: "test/1", Transport: transport})
	resp, err := c.ItemsSearch(context.Background(), "token", api.ItemsSearchRequest{Q: "alpha", Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("ItemsSearch error: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].ID != "i1" {
		t.Fatalf("unexpected response: %+v", resp.Items)
	}
}
