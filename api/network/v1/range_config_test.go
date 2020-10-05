package network

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeConfigValidate(t *testing.T) {

	t.Run("ReturnsNoErrorsWhenRangeConfigIsValid", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.True(t, valid)
		assert.Equal(t, 0, len(errors))
	})

	t.Run("ReturnsErrorWhenTypeIsUnknown", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.Type = TypeUnknown

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
		assert.True(t, strings.HasPrefix(errors[0], "Value for field type is unknown"))
	})

	t.Run("ReturnsErrorWhenRangeTypeIsUnknown", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.RangeType = RangeTypeUnknown

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
		assert.True(t, strings.HasPrefix(errors[0], "Value for field ip_cidr_range_type is unknown"))
	})

	t.Run("ReturnsErrorWhenNetworkIsInvalidCIDR", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "192.0.2.0/388"

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "Value for field network is invalid: invalid CIDR address: 192.0.2.0/388", errors[0])
	})

	t.Run("ReturnsErrorWhenSubnetMaskIsLessThanNetworkIsInvalidCIDR", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "172.28.0.0/14"
		rangeConfig.SubnetMask = 13

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "Value for field subnet_mask is less than the network mask: 13 < 14", errors[0])
	})

	t.Run("ReturnsErrorWhenSubnetMaskIsLessThanZero", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "172.28.0.0/14"
		rangeConfig.SubnetMask = -1

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "Value for field subnet_mask is invalid; it needs to be between 14 and 32", errors[0])
	})

	t.Run("ReturnsErrorWhenSubnetMaskMoreThan32", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "172.28.0.0/14"
		rangeConfig.SubnetMask = -1

		// act
		valid, _, errors := rangeConfig.Validate()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "Value for field subnet_mask is invalid; it needs to be between 14 and 32", errors[0])
	})
}

func TestGetMaxSubnetworkRanges(t *testing.T) {

	t.Run("Returns2IfMaskHasDifferenceOf1", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "172.28.0.0/14"
		rangeConfig.SubnetMask = 15

		// act
		maxSubnetworkRanges := rangeConfig.GetMaxSubnetworkRanges()

		assert.Equal(t, 2, maxSubnetworkRanges)
	})

	t.Run("Returns4IfMaskHasDifferenceOf2", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "172.28.0.0/14"
		rangeConfig.SubnetMask = 16

		// act
		maxSubnetworkRanges := rangeConfig.GetMaxSubnetworkRanges()

		assert.Equal(t, 4, maxSubnetworkRanges)
	})
}

func TestGetAvailableSubnetworkRanges(t *testing.T) {

	t.Run("Returns2RangesIfMaskHasDifferenceOf1", func(t *testing.T) {

		rangeConfig := getValidRangeConfig()
		rangeConfig.NetworkCIDR = "172.28.0.0/14"
		rangeConfig.SubnetMask = 15

		// act
		availableSubnetworkRanges := rangeConfig.GetAvailableSubnetworkRanges()

		assert.Equal(t, 2, len(availableSubnetworkRanges))
		assert.Equal(t, "172.28.0.0/15", availableSubnetworkRanges[0].String())
		assert.Equal(t, "172.30.0.0/15", availableSubnetworkRanges[1].String())
	})
}

func getValidRangeConfig() RangeConfig {
	return RangeConfig{
		Type:        TypeNode,
		RangeType:   RangeTypePrimary,
		NetworkCIDR: "172.28.0.0/14",
		SubnetMask:  21,
	}
}
