package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/codemint/codemint-cli/internal/api"
	"github.com/codemint/codemint-cli/internal/auth"
	"github.com/codemint/codemint-cli/internal/config"
	"github.com/codemint/codemint-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtBy   = "local"
	ctx       appContext
	cfgPath   string
	flagJSON  bool
	flagURL   string
	flagProf  string
	flagDebug bool
)

type appContext struct {
	Config config.Config
	Client *api.Client
	Store  auth.TokenStore
	Mode   output.Mode
}

var rootCmd = &cobra.Command{
	Use:   "codemint",
	Short: "CodeMint CLI",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		mode := output.FromJSONFlag(flagJSON)
		cfg, err := config.Load(config.LoadOptions{ConfigPath: cfgPath, BaseURLOverride: flagURL, ProfileOverride: flagProf})
		if err != nil {
			return err
		}

		store, err := auth.NewTokenStore(cfg.Profile)
		if err != nil {
			return fmt.Errorf("init secure token store: %w", err)
		}

		client := api.NewClient(api.ClientOptions{
			BaseURL:   cfg.BaseURL,
			Timeout:   20 * time.Second,
			UserAgent: fmt.Sprintf("codemint/%s (%s/%s)", version, runtime.GOOS, runtime.GOARCH),
			Debug:     flagDebug,
		})

		ctx = appContext{Config: cfg, Client: client, Store: store, Mode: mode}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitCode := api.ExitCode(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "output JSON")
	rootCmd.PersistentFlags().StringVar(&flagURL, "base-url", "", "override API base URL")
	rootCmd.PersistentFlags().StringVar(&flagProf, "profile", "", "profile name")
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file path")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "enable debug logging")

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newAuthCmd())
	rootCmd.AddCommand(newItemsCmd())
	rootCmd.AddCommand(newOrgCmd())

	rootCmd.AddCommand(newScanCmd())
	rootCmd.AddCommand(newSuggestCmd())
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newSyncCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newRemoveCmd())
	rootCmd.AddCommand(newDoctorCmd())
	rootCmd.AddCommand(newToolCmd())
}

func rootContext() context.Context {
	return context.Background()
}

func tokenFromStore() (string, error) {
	tok, err := ctx.Store.Get(rootContext())
	if err != nil {
		return "", fmt.Errorf("not logged in. run `codemint auth login`: %w", err)
	}
	return tok, nil
}
