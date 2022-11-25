package cms

import (
	"encoding/json"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ContentDir string `json:"content_dir" yaml:"content_dir"`
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
