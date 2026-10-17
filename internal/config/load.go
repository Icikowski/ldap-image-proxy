package config

import "github.com/caarlos0/env/v11"

// Load loads the [Config] from environment variables and validates it.
func Load() (*Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Prefix: "LIP_",
	})
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
