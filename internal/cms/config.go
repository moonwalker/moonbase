package cms

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

const (
	ConfigPath = "moonbase.yaml"
)

type Config struct {
	WorkDir string `json:"workdir" yaml:"workdir"`
}

func ParseConfig(data []byte) *Config {
	cfg := &Config{}

	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		json.Unmarshal(data, cfg)
	}

	return cfg
}
