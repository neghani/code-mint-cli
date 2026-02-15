package catalog

import "testing"

func TestParseRef(t *testing.T) {
	ref, err := ParseRef("@rule/react-best-coding")
	if err != nil {
		t.Fatalf("ParseRef error: %v", err)
	}
	if ref.Type != "rule" || ref.Slug != "react-best-coding" {
		t.Fatalf("unexpected ref: %#v", ref)
	}
}

func TestParseRefInvalid(t *testing.T) {
	_, err := ParseRef("rule/react-best-coding")
	if err == nil {
		t.Fatal("expected error for missing @")
	}
}
