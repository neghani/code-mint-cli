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

func TestAuthMe(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"user":{"id":"u1","email":"dev@example.com","name":"Dev"}}`)),
		}, nil
	})

	c := api.NewClient(api.ClientOptions{BaseURL: "https://example.com", Timeout: 2 * time.Second, UserAgent: "test/1", Transport: transport})
	me, err := c.AuthMe(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("AuthMe error: %v", err)
	}
	if me.Email != "dev@example.com" {
		t.Fatalf("unexpected email: %s", me.Email)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
