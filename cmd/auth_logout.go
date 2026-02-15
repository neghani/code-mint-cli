package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Delete locally stored token",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := ctx.Store.Delete(cmd.Context()); err != nil {
				return err
			}
			fmt.Println("Logged out")
			return nil
		},
	}
}
