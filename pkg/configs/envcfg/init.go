package envcfg

import (
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	stType     string //storage type env in this case
}

func (cfg *Config) postInit() {
	cfg.stType = "env"
}

func (a *Config) GetValue(ctx context.Context) (string, error) {
	return os.Getenv("myVar"), nil
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("Parse environment configuration failed: %w", err)
	}
	cfg.postInit()
	return &cfg, nil
}
