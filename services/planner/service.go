package planner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	networkv1 "github.com/estafette/estafette-gcp-network-planner/api/network/v1"
	"github.com/estafette/estafette-gcp-network-planner/clients/gcp"
	"github.com/rs/zerolog/log"
	computev1 "google.golang.org/api/compute/v1"
)

//go:generate mockgen -package=planner -destination ./mock.go -source=service.go
type Service interface {
	LoadConfig(ctx context.Context) (config *networkv1.Config, err error)
	Suggest(ctx context.Context, region, filter string, networkTypes ...networkv1.Type) (subnetsMap map[networkv1.Type]*net.IPNet, err error)
	SuggestSingleNetworkRange(ctx context.Context, rangeConfigs []networkv1.RangeConfig, subnetworks []*computev1.Subnetwork, region string, networkType networkv1.Type) (subnetworkRange *net.IPNet, err error)
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

	var data []byte

	if s.configPath != "" {
		log.Info().Msgf("Reading config from %v...", s.configPath)
		data, err = ioutil.ReadFile(s.configPath)
		if err != nil {
			return
		}
	} else {
		log.Info().Msg("Reading config from embedded config.json file...")
		data, err = networkv1.Asset("config.json")
		if err != nil {
			return
		}
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	return
}

func (s *service) Suggest(ctx context.Context, region, filter string, networkTypes ...networkv1.Type) (subnetsMap map[networkv1.Type]*net.IPNet, err error) {

	config, err := s.LoadConfig(ctx)
	if err != nil {
		return
	}

	valid, _, errors := config.Validate()
	if !valid {
		return subnetsMap, fmt.Errorf("Config at path %v is not valid: %v", s.configPath, errors)
	}

	projects, err := s.gcpClient.GetProjectByLabels(ctx, []string{filter})
	if err != nil {
		return
	}

	subnetworks, err := s.gcpClient.GetProjectSubnetworks(ctx, projects)
	if err != nil {
		return
	}

	// set default network type
	if len(networkTypes) == 0 {
		networkTypes = []networkv1.Type{
			networkv1.TypeNode,
			networkv1.TypePod,
			networkv1.TypeService,
			networkv1.TypeMaster,
			networkv1.TypeOther,
		}
	}

	// get suggested subnets
	subnetsMap = map[networkv1.Type]*net.IPNet{}
	for _, t := range networkTypes {
		subnetRange, err := s.SuggestSingleNetworkRange(ctx, config.RangeConfigs, subnetworks, region, t)
		if err != nil {
			return subnetsMap, err
		}
		subnetsMap[t] = subnetRange
	}

	// log suggested subnets:
	for k, v := range subnetsMap {
		log.Info().Msgf("%v - %v", k, v)
	}

	return
}

func (s *service) SuggestSingleNetworkRange(ctx context.Context, rangeConfigs []networkv1.RangeConfig, subnetworks []*computev1.Subnetwork, region string, networkType networkv1.Type) (subnetworkRange *net.IPNet, err error) {

	log.Debug().Msgf("Suggesting subnetwork range for region %v and network type %v (with %v range configs and %v subnetworks)...", region, networkType, len(rangeConfigs), len(subnetworks))

	// find range config for region and network type
	filteredRangeConfigs := []networkv1.RangeConfig{}
	for _, rc := range rangeConfigs {
		if rc.Type == networkType && rc.Region == region {
			filteredRangeConfigs = append(filteredRangeConfigs, rc)
		}
	}

	if len(filteredRangeConfigs) == 0 {
		return subnetworkRange, fmt.Errorf("No ranges have been configured for type %v and region %v, can't suggest a subnetwork range", networkType, region)
	}

	if len(filteredRangeConfigs) > 1 {
		return subnetworkRange, fmt.Errorf("Multiple ranges have been configured for type %v and region %v, can't suggest a subnetwork range", networkType, region)
	}

	rangeConfig := filteredRangeConfigs[0]

	// filter subnetworks on whether they're contained in the range config network CIDR
	filteredSubnetworkCIDRs := []string{}
	for _, sn := range subnetworks {

		if !strings.HasSuffix(sn.Region, "/"+rangeConfig.Region) {
			continue
		}

		switch rangeConfig.RangeType {
		case networkv1.RangeTypePrimary:
			contains, err := rangeConfig.ContainsCIDR(sn.IpCidrRange)
			if err != nil {
				return subnetworkRange, err
			}
			if contains {
				filteredSubnetworkCIDRs = append(filteredSubnetworkCIDRs, sn.IpCidrRange)
			}

		case networkv1.RangeTypeSecondary:
			for _, sr := range sn.SecondaryIpRanges {
				contains, err := rangeConfig.ContainsCIDR(sr.IpCidrRange)
				if err != nil {
					return subnetworkRange, err
				}
				if contains {
					filteredSubnetworkCIDRs = append(filteredSubnetworkCIDRs, sn.IpCidrRange)
				}
			}
		}

	}

	log.Debug().Msgf("Filtered subnetworks down to %v applicable subnetworks", len(filteredSubnetworkCIDRs))

	// get first free subnetwork range from rangeconfig
	availableSubnetworkRanges := rangeConfig.GetAvailableSubnetworkRanges()
	for _, subnetRange := range availableSubnetworkRanges {
		// check if it's in use
		rangeIsInUse := false
		for _, snCIDR := range filteredSubnetworkCIDRs {
			subnetworkIP, _, err := net.ParseCIDR(snCIDR)
			if err != nil {
				return subnetworkRange, err
			}

			if subnetRange.Contains(subnetworkIP) {
				log.Debug().Msgf("Range %v is already used by subnet with cidr %v", subnetRange, snCIDR)
				rangeIsInUse = true
				break
			}
		}

		if !rangeIsInUse {
			return subnetRange, nil
		}
	}

	return subnetworkRange, fmt.Errorf("All of the possible %v subnets of range %v are already in use", len(availableSubnetworkRanges), rangeConfig.NetworkCIDR)
}
