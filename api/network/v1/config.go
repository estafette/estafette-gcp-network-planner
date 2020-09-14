package network

type Config struct {
	RangeConfigs []RangeConfig `json:"range_configs"`
}

func (c *Config) Validate() (valid bool, warnings []string, errors []string) {

	// validate all range configs
	for _, rc := range c.RangeConfigs {
		_, w, e := rc.Validate()

		warnings = append(warnings, w...)
		errors = append(errors, e...)
	}

	return len(errors) == 0, warnings, errors
}
