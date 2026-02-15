package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codemint/codemint-cli/internal/api"
	"github.com/codemint/codemint-cli/internal/tooling"
	"github.com/codemint/codemint-cli/internal/util"
)

type InstallResult struct {
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
}

type Manager struct {
	Root string
}

func NewManager(root string) *Manager {
	return &Manager{Root: root}
}

func (m *Manager) ItemDir(tool, itemType string) string {
	switch tool {
	case tooling.ToolCursor:
		if itemType == "skill" {
			return filepath.Join(m.Root, ".cursor", "skills")
		}
		return filepath.Join(m.Root, ".cursor", "rules")
	case tooling.ToolCline:
		if itemType == "skill" {
			return filepath.Join(m.Root, ".cline", "skills")
		}
		return filepath.Join(m.Root, ".clinerules")
	case tooling.ToolWindsurf:
		if itemType == "skill" {
			return filepath.Join(m.Root, ".windsurf", "skills")
		}
		return filepath.Join(m.Root, ".windsurf", "rules")
	case tooling.ToolContinue:
		if itemType == "skill" {
			return filepath.Join(m.Root, ".continue", "skills")
		}
		return filepath.Join(m.Root, ".continue", "rules")
	case tooling.ToolCopilot:
		return filepath.Join(m.Root, ".github", "instructions")
	case tooling.ToolClaude:
		if itemType == "skill" {
			return filepath.Join(m.Root, ".claude", "skills")
		}
		return filepath.Join(m.Root, ".claude", "rules")
	case tooling.ToolCodex:
		return filepath.Join(m.Root, ".codex", itemType+"s")
	default:
		return filepath.Join(m.Root, ".codemint", itemType+"s")
	}
}

func (m *Manager) ItemPath(tool, itemType, slug string) string {
	ext := ".md"
	name := slug
	switch tool {
	case tooling.ToolCursor:
		if itemType == "skill" {
			return filepath.Join(m.ItemDir(tool, itemType), slug, "SKILL.md")
		}
		ext = ".mdc"
	case tooling.ToolCopilot:
		ext = ".instructions.md"
		if itemType == "skill" {
			name = "skill-" + slug
		}
	case tooling.ToolCline:
		if itemType == "skill" {
			return filepath.Join(m.ItemDir(tool, itemType), slug, "SKILL.md")
		}
	case tooling.ToolWindsurf, tooling.ToolContinue, tooling.ToolClaude:
		if itemType == "skill" {
			name = "skill-" + slug
		}
	case tooling.ToolCodex:
		ext = ".md"
	}
	return filepath.Join(m.ItemDir(tool, itemType), name+ext)
}

func (m *Manager) BackupPath(tool, itemType, slug string) string {
	return filepath.Join(m.Root, ".codemint", "backup", itemType+"s", slug+".bak")
}

func (m *Manager) Install(item api.CatalogItem, tool string) (InstallResult, error) {
	if item.Type != "rule" && item.Type != "skill" {
		return InstallResult{}, fmt.Errorf("unsupported item type: %s", item.Type)
	}
	if tool == "" {
		tool = tooling.ToolCodeMint
	}
	content := item.Content
	if content == "" {
		content = defaultContent(item)
	}
	content = renderForTool(tool, item, content)
	path := m.ItemPath(tool, item.Type, item.Slug)
	if err := util.EnsureDir(filepath.Dir(path)); err != nil {
		return InstallResult{}, err
	}
	if existing, err := os.ReadFile(path); err == nil && len(existing) > 0 {
		backup := m.BackupPath(tool, item.Type, item.Slug)
		if err := util.AtomicWriteFile(backup, existing, 0o644); err != nil {
			return InstallResult{}, err
		}
	}
	if err := util.AtomicWriteFile(path, []byte(content), 0o644); err != nil {
		return InstallResult{}, err
	}
	return InstallResult{Path: path, Checksum: util.SHA256Hex([]byte(content))}, nil
}

func (m *Manager) RemovePath(path string) (string, error) {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if empty, _ := util.IsEmptyDir(dir); empty {
			_ = os.Remove(dir)
		}
	}
	return path, nil
}

func defaultContent(item api.CatalogItem) string {
	return fmt.Sprintf("# %s\n\n- Type: %s\n- Ref: @%s/%s\n- Catalog ID: %s\n- Version: %s\n", item.Name, item.Type, item.Type, item.Slug, item.CatalogID, item.Version)
}

func renderForTool(tool string, item api.CatalogItem, content string) string {
	if tool != tooling.ToolCursor {
		return content
	}
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		return content
	}
	return fmt.Sprintf("---\ndescription: %s\nalwaysApply: false\n---\n\n%s\n", safeTitle(item), content)
}

func safeTitle(item api.CatalogItem) string {
	if item.Title != "" {
		return item.Title
	}
	if item.Name != "" {
		return item.Name
	}
	return item.Slug
}
