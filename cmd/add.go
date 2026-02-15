package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/codemint/codemint-cli/internal/catalog"
	"github.com/codemint/codemint-cli/internal/install"
	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var dryRun bool
	var selectedTool string
	cmd := &cobra.Command{
		Use:   "add @rule/<slug>|@skill/<slug>",
		Short: "Install a rule or skill from catalog",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("add expects exactly one identifier")
			}
			ref, err := catalog.ParseRef(args[0])
			if err != nil {
				return err
			}
			tok, err := tokenFromStore()
			if err != nil {
				return err
			}
			item, err := ctx.Client.CatalogGetByRef(c.Context(), tok, ref.Type, ref.Slug)
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
			tool, err := resolveAITool(store, selectedTool, ctx.Mode == output.ModeJSON)
			if err != nil {
				return err
			}
			if idx, ok := manifest.FindByCatalogID(mf.Installed, item.CatalogID); ok {
				ex := mf.Installed[idx]
				if ex.Version == item.Version {
					msg := map[string]any{"status": "unchanged", "ref": ex.Ref, "version": ex.Version, "tool": ex.Tool}
					if ctx.Mode == output.ModeJSON {
						return output.PrintJSON(msg)
					}
					fmt.Printf("Already installed %s@%s\n", ex.Ref, ex.Version)
					return nil
				}
			}
			if dryRun {
				plan := map[string]any{"action": "install", "ref": ref.Raw, "catalogId": item.CatalogID, "version": item.Version, "tool": tool}
				if ctx.Mode == output.ModeJSON {
					return output.PrintJSON(plan)
				}
				fmt.Printf("Dry run: install %s (%s@%s) for %s\n", ref.Raw, item.CatalogID, item.Version, tool)
				return nil
			}

			mgr := install.NewManager(wd)
			installed, err := mgr.Install(*item, tool)
			if err != nil {
				return err
			}
			entry := manifest.Item{
				CatalogID:   item.CatalogID,
				Ref:         ref.Raw,
				Type:        item.Type,
				Slug:        item.Slug,
				Tool:        tool,
				Version:     item.Version,
				Checksum:    firstNonEmpty(item.Checksum, installed.Checksum),
				InstalledAt: time.Now().UTC(),
				Path:        installed.Path,
			}
			if idx, ok := manifest.FindByCatalogID(mf.Installed, entry.CatalogID); ok {
				mf.Installed[idx] = entry
			} else {
				mf.Installed = append(mf.Installed, entry)
			}
			if err := store.Save(mf); err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(entry)
			}
			fmt.Printf("Installed %s (%s)\n", entry.Ref, entry.Version)
			fmt.Printf("Tool: %s\n", entry.Tool)
			fmt.Printf("Path: %s\n", entry.Path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview install without writing files")
	cmd.Flags().StringVar(&selectedTool, "tool", "", "AI coding tool for install target")
	return cmd
}

func firstNonEmpty(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}
