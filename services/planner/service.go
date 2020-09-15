package planner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	networkv1 "github.com/estafette/estafette-gcp-network-planner/api/network/v1"
	"github.com/estafette/estafette-gcp-network-planner/clients/gcp"
	"github.com/rs/zerolog/log"
	computev1 "google.golang.org/api/compute/v1"
)

//go:generate mockgen -package=planner -destination ./mock.go -source=service.go
type Service interface {
	LoadConfig(ctx context.Context) (config *networkv1.Config, err error)
	Suggest(ctx context.Context, region, filter string, networkTypes ...networkv1.Type) (subnetsMap map[networkv1.Type]*net.IPNet, err error)
	SuggestSingleNetworkRange(ctx context.Context, rangeConfigs []networkv1.RangeConfig, subnetworks []*computev1.Subnetwork, routes []*computev1.Route, region string, networkType networkv1.Type) (subnetworkRange *net.IPNet, err error)
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

	routes, err := s.gcpClient.GetProjectRoutes(ctx, projects)
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
		subnetRange, err := s.SuggestSingleNetworkRange(ctx, config.RangeConfigs, subnetworks, routes, region, t)
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

func (s *service) SuggestSingleNetworkRange(ctx context.Context, rangeConfigs []networkv1.RangeConfig, subnetworks []*computev1.Subnetwork, routes []*computev1.Route, region string, networkType networkv1.Type) (subnetworkRange *net.IPNet, err error) {

	log.Debug().Msgf("Suggesting subnetwork range for region %v and network type %v (with %v range configs and %v subnetworks and %v routes)...", region, networkType, len(rangeConfigs), len(subnetworks), len(routes))

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
			overlap, overlapErr := s.rangesOverlap(rangeConfig.NetworkCIDR, sn.IpCidrRange)
			if overlapErr != nil {
				return nil, overlapErr
			}
			if overlap {
				filteredSubnetworkCIDRs = append(filteredSubnetworkCIDRs, sn.IpCidrRange)
			}

		case networkv1.RangeTypeSecondary:
			for _, sr := range sn.SecondaryIpRanges {
				overlap, overlapErr := s.rangesOverlap(rangeConfig.NetworkCIDR, sr.IpCidrRange)
				if overlapErr != nil {
					return nil, overlapErr
				}
				if overlap {
					filteredSubnetworkCIDRs = append(filteredSubnetworkCIDRs, sr.IpCidrRange)
				}
			}
		}
	}
	log.Debug().Msgf("Filtered subnetworks down to %v applicable subnetworks", len(filteredSubnetworkCIDRs))

	// filter routes on whether they're contained in the range config network CIDR
	filteredRouteCIDRs := []string{}
	for _, r := range routes {
		if r.DestRange == "0.0.0.0/0" {
			continue
		}

		overlap, overlapErr := s.rangesOverlap(rangeConfig.NetworkCIDR, r.DestRange)
		if overlapErr != nil {
			return nil, overlapErr
		}
		if overlap {
			filteredRouteCIDRs = append(filteredRouteCIDRs, r.DestRange)
		}
	}
	log.Debug().Msgf("Filtered routes down to %v applicable routes", len(filteredRouteCIDRs))

	// get first free subnetwork range from rangeconfig
	availableSubnetworkRanges := rangeConfig.GetAvailableSubnetworkRanges()
	for i, subnetRange := range availableSubnetworkRanges {
		// check if it's in use by any of the filtered subnets
		rangeIsInUse := false
		for _, snCIDR := range filteredSubnetworkCIDRs {
			overlap, overlapErr := s.rangesOverlap(subnetRange.String(), snCIDR)
			if overlapErr != nil {
				return nil, overlapErr
			}
			if overlap {
				log.Debug().Msgf("Range %v is already used by subnet with cidr %v", subnetRange, snCIDR)
				rangeIsInUse = true
				break
			}
		}
		if !rangeIsInUse {
			// check if it's in use by any of the filtered routes
			for _, rCIDR := range filteredRouteCIDRs {
				overlap, overlapErr := s.rangesOverlap(subnetRange.String(), rCIDR)
				if overlapErr != nil {
					return nil, overlapErr
				}
				if overlap {
					log.Debug().Msgf("Range %v is already used by route with cidr %v", subnetRange, rCIDR)
					rangeIsInUse = true
					break
				}
			}
		}

		if !rangeIsInUse {
			log.Debug().Msgf("%vth range %v of total range %v is available, suggesting it", i, subnetRange, rangeConfig.NetworkCIDR)
			return subnetRange, nil
		}
	}

	return subnetworkRange, fmt.Errorf("All of the possible %v subnets of range %v are already in use", len(availableSubnetworkRanges), rangeConfig.NetworkCIDR)
}

func (s *service) rangesOverlap(cidrA, cidrB string) (overlap bool, err error) {

	_, ipnetA, err := net.ParseCIDR(cidrA)
	if err != nil {
		return false, fmt.Errorf("Parsing %v failed", cidrA)
	}
	firstA, lastA := cidr.AddressRange(ipnetA)

	_, ipnetB, err := net.ParseCIDR(cidrB)
	if err != nil {
		return false, fmt.Errorf("Parsing %v failed", cidrB)
	}
	firstB, lastB := cidr.AddressRange(ipnetB)

	return ipnetA.Contains(firstB) ||
		ipnetA.Contains(lastB) ||
		ipnetB.Contains(firstA) ||
		ipnetB.Contains(lastA), nil
}
