package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {

	t.Run("ReturnsNoErrorsWhenAlleRangeConfigsAreValid", func(t *testing.T) {

		config := getValidConfig()

		// act
		valid, _, errors := config.Validate()

		assert.True(t, valid)
		assert.Equal(t, 0, len(errors))
	})
}

func getValidConfig() Config {
	return Config{
		RangeConfigs: []RangeConfig{
			{
				Type:        TypeNode,
				Region:      "europe-west1",
				RangeType:   RangeTypePrimary,
				NetworkCIDR: "172.28.0.0/14",
				SubnetMask:  21,
			},
			{
				Type:        TypePod,
				Region:      "europe-west1",
				RangeType:   RangeTypeSecondary,
				NetworkCIDR: "10.128.0.0/9",
				SubnetMask:  14,
			},
		},
	}
}
