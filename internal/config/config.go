package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultBaseURL = "https://codemint.app"

type Config struct {
	BaseURL string `json:"base_url"`
	Profile string `json:"profile"`
}

type LoadOptions struct {
	ConfigPath      string
	BaseURLOverride string
	ProfileOverride string
}

func Load(opts LoadOptions) (Config, error) {
	cfg := Config{BaseURL: defaultBaseURL, Profile: "default"}
	path := opts.ConfigPath
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return Config{}, err
		}
		path = filepath.Join(home, ".config", "codemint", "config.json")
	}

	if b, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(b, &cfg); err != nil {
			return Config{}, errors.New("invalid config file")
		}
	}

	if v := EnvBaseURL(); v != "" {
		cfg.BaseURL = v
	}
	if v := EnvProfile(); v != "" {
		cfg.Profile = v
	}
	if opts.BaseURLOverride != "" {
		cfg.BaseURL = opts.BaseURLOverride
	}
	if opts.ProfileOverride != "" {
		cfg.Profile = opts.ProfileOverride
	}
	return cfg, nil
}
