package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/codemint/codemint-cli/internal/api"
	"github.com/codemint/codemint-cli/internal/catalog"
	"github.com/codemint/codemint-cli/internal/install"
	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

type syncPlan struct {
	Upgrade []manifest.Item `json:"upgrade"`
	Same    []manifest.Item `json:"unchanged"`
	Removed []manifest.Item `json:"removed"`
}

func newSyncCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync installed rules/skills with latest catalog versions",
		RunE: func(c *cobra.Command, _ []string) error {
			tok, err := tokenFromStore()
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
			if len(mf.Installed) == 0 {
				if ctx.Mode == output.ModeJSON {
					return output.PrintJSON(syncPlan{})
				}
				fmt.Println("No installed catalog items to sync")
				return nil
			}
			req := api.CatalogSyncRequest{Items: make([]api.CatalogSyncItem, 0, len(mf.Installed))}
			for _, it := range mf.Installed {
				req.Items = append(req.Items, api.CatalogSyncItem{CatalogID: it.CatalogID, Version: it.Version, Checksum: it.Checksum})
			}
			resp, err := ctx.Client.CatalogSync(c.Context(), tok, req)
			if err != nil {
				return err
			}

			mgr := install.NewManager(wd)
			plan := syncPlan{}
			settings, _ := store.LoadSettings()
			for _, local := range mf.Installed {
				result := lookupSync(local.CatalogID, resp.Results)
				if result == nil || result.Removed {
					plan.Removed = append(plan.Removed, local)
					continue
				}
				if result.LatestVersion == "" || result.LatestVersion == local.Version {
					plan.Same = append(plan.Same, local)
					continue
				}
				up := local
				up.Version = result.LatestVersion
				if result.LatestItem.Checksum != "" {
					up.Checksum = result.LatestItem.Checksum
				}
				plan.Upgrade = append(plan.Upgrade, up)
				if dryRun {
					continue
				}
				tool := local.Tool
				if tool == "" {
					tool = settings.AITool
				}
				installed, err := mgr.Install(result.LatestItem, tool)
				if err != nil {
					return err
				}
				up.Path = installed.Path
				up.Tool = tool
				if up.Checksum == "" {
					up.Checksum = installed.Checksum
				}
				up.Ref = catalog.NormalizeRef(local.Type, local.Slug)
				up.InstalledAt = time.Now().UTC()
				if idx, ok := manifest.FindByCatalogID(mf.Installed, local.CatalogID); ok {
					mf.Installed[idx] = up
				}
			}
			if !dryRun {
				if err := store.Save(mf); err != nil {
					return err
				}
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(plan)
			}
			fmt.Printf("Upgrades: %d\n", len(plan.Upgrade))
			fmt.Printf("Unchanged: %d\n", len(plan.Same))
			fmt.Printf("Removed/Deprecated: %d\n", len(plan.Removed))
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview sync plan without writing files")
	return cmd
}

func lookupSync(catalogID string, list []api.CatalogSyncResult) *api.CatalogSyncResult {
	for i := range list {
		if list[i].CatalogID == catalogID {
			return &list[i]
		}
	}
	return nil
}
