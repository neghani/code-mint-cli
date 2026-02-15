package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/codemint/codemint-cli/internal/output"
	"github.com/codemint/codemint-cli/internal/scan"
	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan repository and detect technologies",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("scan accepts at most one path argument")
			}
			target := "."
			if len(args) == 1 {
				target = args[0]
			}
			res, err := scan.Detect(target)
			if err != nil {
				return err
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(res)
			}
			keys := make([]string, 0, len(res.Confidence))
			for k := range res.Confidence {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			rows := make([][]string, 0, len(keys))
			for _, k := range keys {
				rows = append(rows, []string{k, strconv.FormatFloat(res.Confidence[k], 'f', 2, 64)})
			}
			if err := output.PrintTable([]string{"Technology", "Confidence"}, rows); err != nil {
				return err
			}
			return output.PrintTable([]string{"Tag"}, asSingleCol(res.Tags))
		},
	}
	return cmd
}

func asSingleCol(values []string) [][]string {
	rows := make([][]string, 0, len(values))
	for _, v := range values {
		rows = append(rows, []string{v})
	}
	return rows
}
