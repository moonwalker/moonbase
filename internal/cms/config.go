package cms

import (
	"encoding/json"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ContentDir string `json:"contents" yaml:"contents"`
	Components struct {
		Entry        string   `json:"entry" yaml:"entry"`
		Dependencies []string `json:"dependencies" yaml:"dependencies"`
	} `json:"components" yaml:"components"`
}

func ParseConfig(path string, data []byte) *Config {
	cfg := &Config{}
	switch filepath.Ext(path) {
	case ".yaml":
		yaml.Unmarshal(data, cfg)
	case ".json":
		json.Unmarshal(data, cfg)
	}
	return cfg
}
