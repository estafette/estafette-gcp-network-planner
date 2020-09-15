package network

import (
	"fmt"
	"math"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

type RangeConfig struct {
	Type        Type      `json:"type"`
	Region      string    `json:"region"`
	RangeType   RangeType `json:"ip_cidr_range_type"`
	NetworkCIDR string    `json:"network"`
	Comment     string    `json:"comment"`
	SubnetMask  int       `json:"subnet_mask"`
}

func (rc *RangeConfig) Validate() (valid bool, warnings []string, errors []string) {
	// validate type
	if rc.Type == TypeUnknown {
		errors = append(errors, "Value for field type is unknown; please set to node, pod, service, master or other")
	}

	// validate ip_cidr_range_type
	if rc.RangeType == RangeTypeUnknown {
		errors = append(errors, "Value for field ip_cidr_range_type is unknown; please set to primary or secondary")
	}

	// validate network
	_, ipnet, err := net.ParseCIDR(rc.NetworkCIDR)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Value for field network is invalid: %v", err.Error()))
	} else {

		// validate subnet_mask
		ones, _ := ipnet.Mask.Size()
		if rc.SubnetMask >= 0 && rc.SubnetMask <= 32 {
			if rc.SubnetMask < ones {
				errors = append(errors, fmt.Sprintf("Value for field subnet_mask is less than the network mask: %v < %v", rc.SubnetMask, ones))
			}
		} else {
			errors = append(errors, fmt.Sprintf("Value for field subnet_mask is invalid; it needs to be between %v and 32", ones))
		}
	}

	return len(errors) == 0, warnings, errors
}

func (rc *RangeConfig) GetMaxSubnetworkRanges() int {
	_, networkIPnet, err := net.ParseCIDR(rc.NetworkCIDR)
	if err != nil {
		return -1
	}

	ones, _ := networkIPnet.Mask.Size()
	onesDiff := rc.SubnetMask - ones

	return int(math.Pow(2, float64(onesDiff)))
}

func (rc *RangeConfig) GetAvailableSubnetworkRanges() (subnetRanges []*net.IPNet) {

	_, ipnet, err := net.ParseCIDR(rc.NetworkCIDR)
	if err != nil {
		return
	}

	ones, _ := ipnet.Mask.Size()

	i := 0
	maxRanges := rc.GetMaxSubnetworkRanges()
	for i < maxRanges {

		subnetRange, err := cidr.Subnet(ipnet, rc.SubnetMask-ones, i)
		if err == nil {
			subnetRanges = append(subnetRanges, subnetRange)
		}

		i++
	}

	return
}
