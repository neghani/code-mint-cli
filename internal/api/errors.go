package api

import "fmt"

type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string {
	if e.Code == "" {
		return fmt.Sprintf("api error (%d): %s", e.Status, e.Message)
	}
	return fmt.Sprintf("api error (%d/%s): %s", e.Status, e.Code, e.Message)
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	ae, ok := err.(*APIError)
	if !ok {
		return 1
	}
	switch ae.Status {
	case 401:
		return 10
	case 403:
		return 11
	case 422:
		return 12
	case 429:
		return 13
	default:
		if ae.Status >= 500 {
			return 14
		}
		return 1
	}
}
