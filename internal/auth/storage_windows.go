//go:build windows

package auth

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type windowsStore struct {
	target string
}

func newPlatformStore(profile string) (TokenStore, error) {
	return &windowsStore{target: "codemint-" + profile}, nil
}

func (w *windowsStore) Set(_ context.Context, token string) error {
	script := fmt.Sprintf(`cmdkey /generic:%s /user:codemint /pass:%s`, w.target, token)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).CombinedOutput()
	if err != nil {
		return fmt.Errorf("write windows credential manager token: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func (w *windowsStore) Get(_ context.Context) (string, error) {
	return "", errors.New("windows secure credential read requires integration with your org credential tooling")
}

func (w *windowsStore) Delete(_ context.Context) error {
	script := fmt.Sprintf(`cmdkey /delete:%s`, w.target)
	_, _ = exec.Command("powershell", "-NoProfile", "-Command", script).CombinedOutput()
	return nil
}
