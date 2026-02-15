package auth

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/codemint/codemint-cli/internal/api"
)

type LoginOptions struct {
	BaseURL string
	Client  *api.Client
	Store   TokenStore
}

type LoginResult struct {
	Email string
}

var openBrowser = OpenBrowser

func Login(ctx context.Context, opts LoginOptions) (*LoginResult, error) {
	srv, err := newCallbackServer()
	if err != nil {
		return nil, fmt.Errorf("start callback server: %w", err)
	}
	defer func() {
		_ = srv.Close(context.Background())
	}()

	loginURL := fmt.Sprintf("%s/cli-auth?port=%d", trimBaseURL(opts.BaseURL), srv.Port())
	if err := openBrowser(loginURL); err != nil {
		fmt.Printf("Could not open browser automatically. Open this URL manually:\n%s\n", loginURL)
	}

	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	payload, err := srv.WaitForToken(waitCtx)
	if err != nil {
		return nil, fmt.Errorf("waiting for auth callback: %w", err)
	}
	if err := opts.Store.Set(ctx, payload.Token); err != nil {
		return nil, fmt.Errorf("persist token: %w", err)
	}

	me, err := opts.Client.AuthMe(ctx, payload.Token)
	if err != nil {
		_ = opts.Store.Delete(ctx)
		return nil, fmt.Errorf("token verification failed. run `codemint auth login` again: %w", err)
	}

	return &LoginResult{Email: me.Email}, nil
}

func trimBaseURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.String() == "" {
		return raw
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}
