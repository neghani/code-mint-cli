package integration

import (
	"testing"

	"github.com/codemint/codemint-cli/internal/auth"
)

func TestRedactToken(t *testing.T) {
	in := "request failed with Authorization: Bearer abc.def.ghi"
	out := auth.RedactToken(in)
	if out == in {
		t.Fatalf("expected token to be redacted")
	}
}
