package planner

import (
	"context"
	"testing"

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

		projects := []*crmv1.Project{}

		gcpClientMock.
			EXPECT().
			GetProjectByLabels(gomock.Any(), gomock.Any()).
			Return(projects, nil)

		gcpClientMock.
			EXPECT().
			GetProjectSubnetworks(gomock.Any(), gomock.Eq(projects)).
			Return([]*computev1.Subnetwork{}, nil)

		// act
		_, err = service.Suggest(ctx, filter)

		assert.Nil(t, err)
	})
}
