package cmd

import "github.com/spf13/cobra"

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{Use: "auth", Short: "Authentication commands"}
	authCmd.AddCommand(newAuthLoginCmd(), newAuthWhoamiCmd(), newAuthLogoutCmd())
	return authCmd
}
