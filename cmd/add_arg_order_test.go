package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestAddFlagAfterIdentifierViaExecute(t *testing.T) {
	c := newAddCmd()

	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"codemint-add", "foo", "--tool", "cursor"}

	err := c.Execute()
	if err == nil {
		t.Fatalf("expected error for invalid identifier")
	}
	if strings.Contains(err.Error(), "add expects exactly one identifier") {
		t.Fatalf("flag-after-arg parsing failed: %v", err)
	}
	if !strings.Contains(err.Error(), "invalid identifier") {
		t.Fatalf("unexpected error: %v", err)
	}
}
