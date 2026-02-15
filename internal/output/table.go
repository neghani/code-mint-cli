package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintTable(headers []string, rows [][]string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, strings.Join(headers, "\t")); err != nil {
		return err
	}
	for _, r := range rows {
		if _, err := fmt.Fprintln(w, strings.Join(r, "\t")); err != nil {
			return err
		}
	}
	return w.Flush()
}
