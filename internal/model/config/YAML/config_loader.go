package yaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*RouteConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config RouteConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
