package cmd

import (
	"fmt"
	"os"
	"runtime"

	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

// rootCmd represents the base command when called without any subcommands
var (
	verbose     bool
	concurrency int

	rootCmd = &cobra.Command{
		Use:   "gcp-network-planner",
		Short: "The command-line interface for planning GCP networks",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	foundation.InitLoggingByFormatSilent(foundation.NewApplicationInfo(appgroup, app, version, branch, revision, buildDate), foundation.LogFormatConsole)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 5, "level of concurrency")
}
