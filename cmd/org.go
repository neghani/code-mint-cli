package cmd

import "github.com/spf13/cobra"

func newOrgCmd() *cobra.Command {
	org := &cobra.Command{Use: "org", Short: "Organization commands"}
	org.AddCommand(newOrgListCmd())
	return org
}
