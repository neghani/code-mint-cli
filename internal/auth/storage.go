package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

type TokenStore interface {
	Set(ctx context.Context, token string) error
	Get(ctx context.Context) (string, error)
	Delete(ctx context.Context) error
}

func NewTokenStore(profile string) (TokenStore, error) {
	if profile == "" {
		profile = "default"
	}
	return newPlatformStore(profile)
}

type fileStore struct {
	path string
}

func newFileStore(profile string) (*fileStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".config", "codemint")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &fileStore{path: filepath.Join(dir, "token-"+profile)}, nil
}

func (f *fileStore) Set(_ context.Context, token string) error {
	return os.WriteFile(f.path, []byte(token), 0o600)
}

func (f *fileStore) Get(_ context.Context) (string, error) {
	b, err := os.ReadFile(f.path)
	if err != nil {
		return "", errors.New("no local token found")
	}
	return string(b), nil
}

func (f *fileStore) Delete(_ context.Context) error {
	if err := os.Remove(f.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
