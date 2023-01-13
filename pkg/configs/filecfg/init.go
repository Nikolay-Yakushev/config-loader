package filecfg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"
)

type Config struct {
	configPath string
	myVar      string
	stType     string //storage type env in this case
}

func getConfigPath() (string, error) {
	curPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error occured while gettinmg working directory: %w", err)
	}
	cfgPath := fmt.Sprintf("%s/config.yaml", filepath.Dir(curPath))
	return cfgPath, nil
}

func (cfg *Config) postInit() error {
	cfg.stType = "file"
	cfgPath, err := getConfigPath()
	if err != nil {
		err = fmt.Errorf("error occured %w", err)
	}
	cfg.configPath = cfgPath
	return err
}

func (a *Config) GetValue(ctx context.Context) (string, error) {
	var (
		value string
		confFile map[string]string
	)
	buf, err := ioutil.ReadFile(a.configPath)
	if err != nil {
		return value, fmt.Errorf("error read file: %w", err)
	}
	unMarshalErr := yaml.Unmarshal(buf, &confFile)
	if unMarshalErr != nil {
		return value, fmt.Errorf("Error during parsing data %w", err)
	}
	value, ok := confFile["myVar"]
	if !ok{
		return value, fmt.Errorf("No value found in config %w", err)
	}
	return value, nil

}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("Parse environment configuration failed: %w", err)
	}
	err := cfg.postInit()
	return &cfg, err
}
