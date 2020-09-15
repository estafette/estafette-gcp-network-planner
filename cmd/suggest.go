package cmd

import (
	"github.com/estafette/estafette-gcp-network-planner/clients/gcp"
	"github.com/estafette/estafette-gcp-network-planner/services/planner"
	"github.com/spf13/cobra"
)

var (
	region string
	filter string
)

func init() {
	rootCmd.AddCommand(suggestCmd)

	// command-specific flags
	suggestCmd.Flags().StringVar(&region, "region", "europe-west1", "Region to request subnetwork range for")
	suggestCmd.Flags().StringVar(&filter, "filter", "", "Filter for limiting projects to retrieve existing network ranges for, see https://cloud.google.com/resource-manager/reference/rest/v1/projects/list#query-parameters")
}

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Suggest a free network range for a subnetwork",
	RunE: func(cmd *cobra.Command, args []string) error {

		// init gcp client
		gcpClient, err := gcp.NewClient(cmd.Context(), concurrency)
		if err != nil {
			return err
		}

		// init planner service
		plannerService, err := planner.NewService(cmd.Context(), gcpClient, configFilePath)
		if err != nil {
			return err
		}

		_, err = plannerService.Suggest(cmd.Context(), region, filter)
		if err != nil {
			return err
		}

		return nil
	},
}
