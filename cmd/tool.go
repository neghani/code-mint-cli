package cmd

import (
	"fmt"
	"os"

	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/codemint/codemint-cli/internal/tooling"
	"github.com/spf13/cobra"
)

func newToolCmd() *cobra.Command {
	toolCmd := &cobra.Command{
		Use:   "tool",
		Short: "Manage default AI coding tool for this repository",
	}
	toolCmd.AddCommand(newToolSetCmd(), newToolCurrentCmd(), newToolListCmd())
	return toolCmd
}

func newToolSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <tool>",
		Short: "Set default AI coding tool",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("tool set expects exactly one tool name")
			}
			tool := args[0]
			if err := tooling.Validate(tool); err != nil {
				return err
			}
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			store := manifest.New(wd)
			settings, err := store.LoadSettings()
			if err != nil {
				return err
			}
			settings.AITool = tool
			if err := store.SaveSettings(settings); err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(map[string]any{"tool": tool, "saved": true})
			}
			fmt.Printf("Default AI tool set to %s\n", tool)
			return nil
		},
	}
}

func newToolCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show default AI coding tool",
		RunE: func(_ *cobra.Command, _ []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			settings, err := manifest.New(wd).LoadSettings()
			if err != nil {
				return err
			}
			if settings.AITool == "" {
				if ctx.Mode == output.ModeJSON {
					return output.PrintJSON(map[string]any{"tool": "", "configured": false})
				}
				fmt.Println("No default tool set")
				return nil
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(map[string]any{"tool": settings.AITool, "configured": true})
			}
			fmt.Println(settings.AITool)
			return nil
		},
	}
}

func newToolListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List supported AI coding tools",
		RunE: func(_ *cobra.Command, _ []string) error {
			tools := tooling.Supported()
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(map[string]any{"tools": tools})
			}
			rows := make([][]string, 0, len(tools))
			for _, t := range tools {
				rows = append(rows, []string{t})
			}
			return output.PrintTable([]string{"Tool"}, rows)
		},
	}
}
