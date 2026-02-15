package cmd

import (
	"strconv"
	"strings"

	"github.com/codemint/codemint-cli/internal/api"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

func newItemsSearchCmd() *cobra.Command {
	var q string
	var itemType string
	var tags []string
	var page int
	var limit int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search items",
		RunE: func(c *cobra.Command, _ []string) error {
			tok, err := tokenFromStore()
			if err != nil {
				return err
			}
			resp, err := ctx.Client.ItemsSearch(c.Context(), tok, api.ItemsSearchRequest{
				Q:     q,
				Type:  itemType,
				Tags:  tags,
				Page:  page,
				Limit: limit,
			})
			if err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(resp)
			}

			rows := make([][]string, 0, len(resp.Items))
			for _, it := range resp.Items {
				rows = append(rows, []string{it.ID, it.Name, it.Type, strings.Join(it.Tags, ","), strconv.Itoa(it.Score)})
			}
			return output.PrintTable([]string{"ID", "Name", "Type", "Tags", "Score"}, rows)
		},
	}

	cmd.Flags().StringVar(&q, "q", "", "query string")
	cmd.Flags().StringVar(&itemType, "type", "", "item type filter")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "comma-separated tags")
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&limit, "limit", 20, "page size")
	_ = cmd.MarkFlagRequired("q")
	return cmd
}
