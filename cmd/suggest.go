package cmd

import (
	"fmt"
	"strings"

	"github.com/codemint/codemint-cli/internal/api"
	"github.com/codemint/codemint-cli/internal/catalog"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/codemint/codemint-cli/internal/scan"
	"github.com/spf13/cobra"
)

type suggestItem struct {
	Ref    string   `json:"ref"`
	Reason string   `json:"reason"`
	Tags   []string `json:"tags"`
}

func newSuggestCmd() *cobra.Command {
	var path string
	var onlyType string

	cmd := &cobra.Command{
		Use:   "suggest",
		Short: "Suggest rules and skills based on repository scan",
		RunE: func(c *cobra.Command, _ []string) error {
			tok, err := tokenFromStore()
			if err != nil {
				return err
			}
			res, err := scan.Detect(path)
			if err != nil {
				return err
			}

			types := []string{catalog.TypeRule, catalog.TypeSkill}
			if onlyType != "" {
				types = []string{onlyType}
			}
			recs := make([]suggestItem, 0)
			for _, t := range types {
				items, err := ctx.Client.CatalogSuggest(c.Context(), tok, api.CatalogLookupRequest{Type: t, Tags: res.Tags, Q: strings.Join(res.Tags, " ")})
				if err != nil {
					return err
				}
				for _, it := range items {
					reason := reasonForItem(it, res.Tags)
					recs = append(recs, suggestItem{Ref: catalog.NormalizeRef(it.Type, it.Slug), Reason: reason, Tags: it.Tags})
				}
			}

			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(map[string]any{"scan": res, "suggestions": recs})
			}
			if len(recs) == 0 {
				fmt.Println("No suggestions found for detected stack")
				return nil
			}
			rows := make([][]string, 0, len(recs))
			for _, r := range recs {
				rows = append(rows, []string{r.Ref, r.Reason})
			}
			return output.PrintTable([]string{"Recommendation", "Reason"}, rows)
		},
	}
	cmd.Flags().StringVar(&path, "path", ".", "path to project")
	cmd.Flags().StringVar(&onlyType, "type", "", "filter by type: rule or skill")
	return cmd
}

func reasonForItem(item api.CatalogItem, tags []string) string {
	matches := 0
	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}
	for _, t := range item.Tags {
		if _, ok := tagSet[t]; ok {
			matches++
		}
	}
	if matches > 0 {
		return fmt.Sprintf("Matched %d detected tags", matches)
	}
	if item.Deprecated {
		return "Deprecated item; avoid unless required"
	}
	return "General recommendation from catalog"
}
