package cmd

import (
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

func newAuthWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show current authenticated user",
		RunE: func(cmd *cobra.Command, _ []string) error {
			tok, err := tokenFromStore()
			if err != nil {
				return err
			}
			me, err := ctx.Client.AuthMe(cmd.Context(), tok)
			if err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(me)
			}
			return output.PrintTable([]string{"ID", "Email", "Name"}, [][]string{{me.ID, me.Email, me.Name}})
		},
	}
}
