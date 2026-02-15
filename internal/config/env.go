package config

import "os"

func EnvBaseURL() string {
	return os.Getenv("CODEMINT_BASE_URL")
}

func EnvProfile() string {
	return os.Getenv("CODEMINT_PROFILE")
}
