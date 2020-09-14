package cmd

import (
	"context"

	foundation "github.com/estafette/estafette-foundation"
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

		// create context to cancel commands on sigterm
		ctx := foundation.InitCancellationContext(context.Background())

		// init gcp client
		gcpClient, err := gcp.NewClient(ctx, concurrency)
		if err != nil {
			return err
		}

		// init planner service
		plannerService, err := planner.NewService(ctx, gcpClient, configFilePath)
		if err != nil {
			return err
		}

		_, err = plannerService.Suggest(ctx, region, filter)
		if err != nil {
			return err
		}

		return nil
	},
}
