package cmd

import (
	"fmt"
	"os"

	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var installedOnly bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List local codemint installs",
		RunE: func(_ *cobra.Command, _ []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			mf, err := manifest.New(wd).Load()
			if err != nil {
				return err
			}
			if !installedOnly {
				if ctx.Mode == output.ModeJSON {
					return output.PrintJSON(mf)
				}
			}
			if len(mf.Installed) == 0 {
				fmt.Println("No installed items")
				return nil
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(mf.Installed)
			}
			rows := make([][]string, 0, len(mf.Installed))
			for _, it := range mf.Installed {
				rows = append(rows, []string{it.Ref, it.Tool, it.Version, it.CatalogID, it.Path})
			}
			return output.PrintTable([]string{"Item", "Tool", "Version", "Catalog ID", "Path"}, rows)
		},
	}
	cmd.Flags().BoolVar(&installedOnly, "installed", false, "show installed items")
	return cmd
}
