package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/codemint/codemint-cli/internal/api"
)

type callbackServer struct {
	server *http.Server
	ln     net.Listener
	tokens chan api.CLIAuthCallbackPayload
}

func newCallbackServer() (*callbackServer, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	cs := &callbackServer{ln: ln, tokens: make(chan api.CLIAuthCallbackPayload, 1)}
	mux := http.NewServeMux()
	mux.HandleFunc("/", cs.handle)
	mux.HandleFunc("/callback", cs.handle)
	cs.server = &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		_ = cs.server.Serve(ln)
	}()
	return cs, nil
}

func (s *callbackServer) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

func (s *callbackServer) WaitForToken(ctx context.Context) (api.CLIAuthCallbackPayload, error) {
	select {
	case <-ctx.Done():
		return api.CLIAuthCallbackPayload{}, ctx.Err()
	case tok := <-s.tokens:
		return tok, nil
	}
}

func (s *callbackServer) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *callbackServer) handle(w http.ResponseWriter, r *http.Request) {
	payload := api.CLIAuthCallbackPayload{Token: r.URL.Query().Get("token"), ExpiresAt: r.URL.Query().Get("expiresAt")}
	if payload.Token == "" && r.Method == http.MethodPost {
		_ = json.NewDecoder(r.Body).Decode(&payload)
	}
	if payload.Token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}
	select {
	case s.tokens <- payload:
	default:
	}
	_, _ = fmt.Fprintln(w, "CodeMint CLI authentication complete. You can close this tab.")
}
