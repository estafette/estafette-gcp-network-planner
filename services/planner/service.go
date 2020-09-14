package planner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	networkv1 "github.com/estafette/estafette-gcp-network-planner/api/network/v1"
	"github.com/estafette/estafette-gcp-network-planner/clients/gcp"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -package=planner -destination ./mock.go -source=service.go
type Service interface {
	LoadConfig(ctx context.Context) (config *networkv1.Config, err error)
	Suggest(ctx context.Context, filter string) (rangeConfigs []networkv1.RangeConfig, err error)
}

func NewService(ctx context.Context, gcpClient gcp.Client, configPath string) (Service, error) {
	return &service{
		gcpClient:  gcpClient,
		configPath: configPath,
	}, nil
}

type service struct {
	gcpClient  gcp.Client
	configPath string
}

func (s *service) LoadConfig(ctx context.Context) (config *networkv1.Config, err error) {

	log.Info().Msgf("Reading config from %v...", s.configPath)

	data, err := ioutil.ReadFile(s.configPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	return
}

func (s *service) Suggest(ctx context.Context, filter string) (rangeConfigs []networkv1.RangeConfig, err error) {

	config, err := s.LoadConfig(ctx)
	if err != nil {
		return
	}

	valid, _, errors := config.Validate()
	if !valid {
		return rangeConfigs, fmt.Errorf("Config at path %v is not valid: %v", s.configPath, errors)
	}

	projects, err := s.gcpClient.GetProjectByLabels(ctx, []string{filter})
	if err != nil {
		return
	}

	subnetworks, err := s.gcpClient.GetProjectSubnetworks(ctx, projects)
	if err != nil {
		return
	}

	log.Info().Interface("subnetworks", subnetworks).Msgf("Retrieved all subnetworks for projects with filter %v", filter)

	return
}
