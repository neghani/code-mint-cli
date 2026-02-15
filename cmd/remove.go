package cmd

import (
	"fmt"
	"os"

	"github.com/codemint/codemint-cli/internal/catalog"
	"github.com/codemint/codemint-cli/internal/install"
	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove @rule/<slug>|@skill/<slug>",
		Short: "Remove installed rule or skill",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("remove expects exactly one identifier")
			}
			ref, err := catalog.ParseRef(args[0])
			if err != nil {
				return err
			}
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			store := manifest.New(wd)
			mf, err := store.Load()
			if err != nil {
				return err
			}
			idx, ok := manifest.FindByRef(mf.Installed, ref.Raw)
			if !ok {
				return fmt.Errorf("not installed: %s", ref.Raw)
			}
			mgr := install.NewManager(wd)
			targetPath := mf.Installed[idx].Path
			if targetPath == "" {
				targetPath = mgr.ItemPath(mf.Installed[idx].Tool, ref.Type, ref.Slug)
			}
			path, err := mgr.RemovePath(targetPath)
			if err != nil {
				return err
			}
			mf.Installed = append(mf.Installed[:idx], mf.Installed[idx+1:]...)
			if err := store.Save(mf); err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(map[string]any{"removed": ref.Raw, "path": path})
			}
			fmt.Printf("Removed %s\n", ref.Raw)
			return nil
		},
	}
	return cmd
}
