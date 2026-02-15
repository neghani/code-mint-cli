package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	dir := t.TempDir()
	files := []string{"package.json", "tsconfig.json", "next.config.ts", "schema.prisma", "Dockerfile"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("{}"), 0o644); err != nil {
			t.Fatalf("write %s: %v", f, err)
		}
	}
	res, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect error: %v", err)
	}
	if len(res.Tags) == 0 {
		t.Fatal("expected non-empty tags")
	}
	if res.Confidence["nextjs"] < 0.9 {
		t.Fatalf("expected nextjs confidence >= 0.9, got %f", res.Confidence["nextjs"])
	}
}
