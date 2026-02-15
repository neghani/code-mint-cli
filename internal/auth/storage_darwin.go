//go:build darwin

package auth

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type darwinStore struct {
	account string
	service string
}

func newPlatformStore(profile string) (TokenStore, error) {
	return &darwinStore{account: profile, service: "codemint-cli"}, nil
}

func (d *darwinStore) Set(_ context.Context, token string) error {
	cmd := exec.Command("security", "add-generic-password", "-a", d.account, "-s", d.service, "-w", token, "-U")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("write macOS keychain token: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func (d *darwinStore) Get(_ context.Context) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", d.account, "-s", d.service, "-w")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("read macOS keychain token: %s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func (d *darwinStore) Delete(_ context.Context) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", d.account, "-s", d.service)
	_, _ = cmd.CombinedOutput()
	return nil
}
