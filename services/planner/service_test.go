package planner

import (
	"context"
	"testing"

	networkv1 "github.com/estafette/estafette-gcp-network-planner/api/network/v1"
	"github.com/estafette/estafette-gcp-network-planner/clients/gcp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	crmv1 "google.golang.org/api/cloudresourcemanager/v1"
	computev1 "google.golang.org/api/compute/v1"
)

func TestSuggest(t *testing.T) {

	t.Run("ReturnsSuggestionsForNodePodAndServiceRanges", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")
		filter := "labels.environment=dev"
		region := "europe-west1"

		projects := []*crmv1.Project{}

		gcpClientMock.
			EXPECT().
			GetProjectByLabels(gomock.Any(), gomock.Any()).
			Return(projects, nil)

		gcpClientMock.
			EXPECT().
			GetProjectSubnetworks(gomock.Any(), gomock.Eq(projects)).
			Return([]*computev1.Subnetwork{}, nil)

		gcpClientMock.
			EXPECT().
			GetProjectRoutes(gomock.Any(), gomock.Eq(projects)).
			Return([]*computev1.Route{}, nil)

		// act
		_, err = service.Suggest(ctx, region, filter)

		assert.Nil(t, err)
	})
}

func TestSuggestSingleNetworkRange(t *testing.T) {

	t.Run("ReturnsErrorWhenNoRangeConfigsMatchRegionAndNetworkType", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{}
		subnetworks := []*computev1.Subnetwork{}
		routes := []*computev1.Route{}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		_, err = service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.NotNil(t, err)
		assert.Equal(t, "No ranges have been configured for type node and region europe-west1, can't suggest a subnetwork range", err.Error())
	})

	t.Run("ReturnsErrorWhenMoreThanOneRangeConfigsMatchRegionAndNetworkType", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  21,
			},
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypeSecondary,
				NetworkCIDR: "10.128.0.0/9",
				SubnetMask:  14,
			},
		}
		subnetworks := []*computev1.Subnetwork{}
		routes := []*computev1.Route{}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		_, err = service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.NotNil(t, err)
		assert.Equal(t, "Multiple ranges have been configured for type node and region europe-west1, can't suggest a subnetwork range", err.Error())
	})

	t.Run("ReturnsErrorWhenMoreOneRangeConfigMatchesAndAllPossibleSubnetsAreInUseBySubnets", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  15,
			},
		}
		subnetworks := []*computev1.Subnetwork{
			{
				IpCidrRange: "172.28.0.0/15",
				Region:      "https://www.googleapis.com/compute/v1/projects/project-id/regions/europe-west1",
			},
			{
				IpCidrRange: "172.30.0.0/15",
				Region:      "https://www.googleapis.com/compute/v1/projects/project-id/regions/europe-west1",
			},
		}
		routes := []*computev1.Route{}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		_, err = service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.NotNil(t, err)
		assert.Equal(t, "All of the possible 2 subnets of range 172.28.0.0/14 are already in use", err.Error())
	})

	t.Run("ReturnsFirstAvailableRangeIfSomeOfThemAreInUseBySubnets", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  15,
			},
		}
		subnetworks := []*computev1.Subnetwork{
			{
				IpCidrRange: "172.28.0.0/15",
				Region:      "https://www.googleapis.com/compute/v1/projects/project-id/regions/europe-west1",
			},
		}
		routes := []*computev1.Route{}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		subnetworkRange, err := service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.Nil(t, err)
		assert.Equal(t, "172.30.0.0/15", subnetworkRange.String())
	})

	t.Run("ReturnsFirstAvailableRangeIfSomeOfThemAreInUseBySecondarySubnetRanges", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypePod,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypeSecondary,
				NetworkCIDR: "10.0.0.0/9",
				SubnetMask:  16,
			},
		}
		subnetworks := []*computev1.Subnetwork{
			{
				IpCidrRange: "172.28.0.0/15",
				Region:      "https://www.googleapis.com/compute/v1/projects/project-id/regions/europe-west1",
				SecondaryIpRanges: []*computev1.SubnetworkSecondaryRange{
					{
						IpCidrRange: "10.0.0.0/16",
					},
				},
			},
		}
		routes := []*computev1.Route{}
		region := "europe-west1"
		networkType := networkv1.TypePod

		// act
		subnetworkRange, err := service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.Nil(t, err)
		assert.Equal(t, "10.1.0.0/16", subnetworkRange.String())
	})

	t.Run("ReturnsErrorWhenMoreOneRangeConfigMatchesAndAllPossibleSubnetsAreInUseByRoutes", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  15,
			},
		}
		subnetworks := []*computev1.Subnetwork{}
		routes := []*computev1.Route{
			{
				DestRange: "172.28.0.0/15",
			},
			{
				DestRange: "172.30.0.0/15",
			},
		}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		_, err = service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.NotNil(t, err)
		assert.Equal(t, "All of the possible 2 subnets of range 172.28.0.0/14 are already in use", err.Error())
	})

	t.Run("ReturnsFirstAvailableRangeIfSomeOfThemAreInUseByRoutes", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  15,
			},
		}
		subnetworks := []*computev1.Subnetwork{}
		routes := []*computev1.Route{
			{
				DestRange: "172.28.0.0/15",
			},
		}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		subnetworkRange, err := service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.Nil(t, err)
		assert.NotNil(t, subnetworkRange)
		assert.Equal(t, "172.30.0.0/15", subnetworkRange.String())
	})

	t.Run("ReturnsFirstRangeIfNoneOfThemAreInUse", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  15,
			},
		}
		subnetworks := []*computev1.Subnetwork{}
		routes := []*computev1.Route{}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		subnetworkRange, err := service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.Nil(t, err)
		assert.NotNil(t, subnetworkRange)
		assert.Equal(t, "172.28.0.0/15", subnetworkRange.String())
	})

	t.Run("ExludeRoutesWithAllZeroes", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gcpClientMock := gcp.NewMockClient(ctrl)

		ctx := context.Background()
		service, err := NewService(ctx, gcpClientMock, "./test-config.json")

		rangeConfigs := []networkv1.RangeConfig{
			{
				Type:        networkv1.TypeNode,
				Region:      "europe-west1",
				RangeType:   networkv1.RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  15,
			},
		}
		subnetworks := []*computev1.Subnetwork{}
		routes := []*computev1.Route{
			{
				DestRange: "0.0.0.0/0",
			},
		}
		region := "europe-west1"
		networkType := networkv1.TypeNode

		// act
		subnetworkRange, err := service.SuggestSingleNetworkRange(ctx, rangeConfigs, subnetworks, routes, region, networkType)

		assert.Nil(t, err)
		assert.NotNil(t, subnetworkRange)
		assert.Equal(t, "172.28.0.0/15", subnetworkRange.String())
	})
}
