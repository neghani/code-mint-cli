package cmd

import (
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

func newOrgListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List organizations for current user",
		RunE: func(cmd *cobra.Command, _ []string) error {
			tok, err := tokenFromStore()
			if err != nil {
				return err
			}
			orgs, err := ctx.Client.OrgList(cmd.Context(), tok)
			if err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(orgs)
			}
			rows := make([][]string, 0, len(orgs.Organizations))
			for _, o := range orgs.Organizations {
				rows = append(rows, []string{o.ID, o.Slug, o.Name, o.Role})
			}
			return output.PrintTable([]string{"ID", "Slug", "Name", "Role"}, rows)
		},
	}
}
