package network

import (
	"fmt"
	"net"
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
		if rc.SubnetMask < ones {
			errors = append(errors, fmt.Sprintf("Value for field subnet_mask is less than the network mask: %v < %v", rc.SubnetMask, ones))
		}
	}

	return len(errors) == 0, warnings, errors
}
