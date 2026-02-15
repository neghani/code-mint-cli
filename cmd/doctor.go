package cmd

import (
	"fmt"
	"os"

	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

type doctorCheck struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

func newDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Validate auth, manifest, and managed paths",
		RunE: func(c *cobra.Command, _ []string) error {
			checks := make([]doctorCheck, 0, 4)

			tok, err := ctx.Store.Get(c.Context())
			if err != nil || tok == "" {
				checks = append(checks, doctorCheck{Name: "auth token", OK: false, Detail: "missing token; run codemint auth login"})
			} else {
				checks = append(checks, doctorCheck{Name: "auth token", OK: true, Detail: "token available in secure store"})
			}

			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			ms := manifest.New(wd)
			mf, err := ms.Load()
			if err != nil {
				checks = append(checks, doctorCheck{Name: "manifest", OK: false, Detail: err.Error()})
			} else {
				checks = append(checks, doctorCheck{Name: "manifest", OK: true, Detail: fmt.Sprintf("%d installed item(s)", len(mf.Installed))})
			}
			settings, err := ms.LoadSettings()
			if err != nil {
				checks = append(checks, doctorCheck{Name: "ai tool", OK: false, Detail: err.Error()})
			} else if settings.AITool == "" {
				checks = append(checks, doctorCheck{Name: "ai tool", OK: false, Detail: "not selected yet; run `codemint tool set <name>` or first `codemint add` will prompt"})
			} else {
				checks = append(checks, doctorCheck{Name: "ai tool", OK: true, Detail: settings.AITool})
			}
			for _, dir := range []string{ms.BaseDir(), wd} {
				st, err := os.Stat(dir)
				if err != nil {
					checks = append(checks, doctorCheck{Name: "path", OK: false, Detail: dir + ": " + err.Error()})
					continue
				}
				checks = append(checks, doctorCheck{Name: "path", OK: st.IsDir(), Detail: dir})
			}
			if ctx.Mode == output.ModeJSON {
				return output.PrintJSON(checks)
			}
			rows := make([][]string, 0, len(checks))
			for _, ch := range checks {
				status := "FAIL"
				if ch.OK {
					status = "OK"
				}
				rows = append(rows, []string{status, ch.Name, ch.Detail})
			}
			return output.PrintTable([]string{"Status", "Check", "Detail"}, rows)
		},
	}
	return cmd
}
