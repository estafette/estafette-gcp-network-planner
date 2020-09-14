package cmd

import (
	"fmt"
	"os"
	"runtime"

	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	verbose        bool
	concurrency    int
	configFilePath string

	rootCmd = &cobra.Command{
		Use:   "gcp-network-planner",
		Short: "The command-line interface for planning GCP networks",
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 5, "level of concurrency")
	rootCmd.PersistentFlags().StringVar(&configFilePath, "config-file", "", "path to config file")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	foundation.InitLoggingByFormatSilent(foundation.NewApplicationInfo(appgroup, app, version, branch, revision, buildDate), foundation.LogFormatConsole)

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Verbose mode enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		// log.Info().Msg("Verbose mode disabled")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
