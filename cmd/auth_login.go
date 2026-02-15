package cmd

import (
	"fmt"

	"github.com/codemint/codemint-cli/internal/auth"
	"github.com/spf13/cobra"
)

func newAuthLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Log in using browser flow",
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := auth.Login(cmd.Context(), auth.LoginOptions{
				BaseURL: ctx.Config.BaseURL,
				Client:  ctx.Client,
				Store:   ctx.Store,
			})
			if err != nil {
				return err
			}
			fmt.Printf("Logged in as %s\n", res.Email)
			return nil
		},
	}
}
