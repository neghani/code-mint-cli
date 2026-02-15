package tooling

import "fmt"

const (
	ToolCursor   = "cursor"
	ToolCline    = "cline"
	ToolWindsurf = "windsurf"
	ToolContinue = "continue"
	ToolCopilot  = "copilot"
	ToolClaude   = "claude"
	ToolCodex    = "codex"
	ToolCodeMint = "codemint"
)

var supported = []string{ToolCursor, ToolCline, ToolWindsurf, ToolContinue, ToolCopilot, ToolClaude, ToolCodex}

func Supported() []string {
	out := make([]string, len(supported))
	copy(out, supported)
	return out
}

func Validate(tool string) error {
	if tool == "" {
		return fmt.Errorf("tool is required")
	}
	for _, t := range supported {
		if tool == t {
			return nil
		}
	}
	return fmt.Errorf("unsupported tool %q (supported: %v)", tool, supported)
}
