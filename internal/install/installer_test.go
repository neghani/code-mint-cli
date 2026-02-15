package install

import (
	"path/filepath"
	"testing"

	"github.com/codemint/codemint-cli/internal/tooling"
)

func TestItemPathByTool(t *testing.T) {
	m := NewManager("/repo")
	cases := []struct {
		tool     string
		itemType string
		slug     string
		expect   string
	}{
		{tooling.ToolCursor, "rule", "react-best", filepath.Join("/repo", ".cursor", "rules", "react-best.mdc")},
		{tooling.ToolCursor, "skill", "node-js", filepath.Join("/repo", ".cursor", "skills", "node-js", "SKILL.md")},
		{tooling.ToolCline, "rule", "safe-api", filepath.Join("/repo", ".clinerules", "safe-api.md")},
		{tooling.ToolCline, "skill", "node-js", filepath.Join("/repo", ".cline", "skills", "node-js", "SKILL.md")},
		{tooling.ToolCopilot, "rule", "secure", filepath.Join("/repo", ".github", "instructions", "secure.instructions.md")},
		{tooling.ToolCodex, "rule", "secure", filepath.Join("/repo", ".codex", "rules", "secure.md")},
		{tooling.ToolCodex, "skill", "node-js", filepath.Join("/repo", ".codex", "skills", "node-js.md")},
	}
	for _, tc := range cases {
		got := m.ItemPath(tc.tool, tc.itemType, tc.slug)
		if got != tc.expect {
			t.Fatalf("ItemPath(%s,%s,%s) got %s want %s", tc.tool, tc.itemType, tc.slug, got, tc.expect)
		}
	}
}
