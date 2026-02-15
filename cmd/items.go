package cmd

import "github.com/spf13/cobra"

func newItemsCmd() *cobra.Command {
	items := &cobra.Command{Use: "items", Short: "Items commands"}
	items.AddCommand(newItemsSearchCmd())
	return items
}
