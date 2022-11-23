package cms

import (
	"encoding/json"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Content struct {
		Dir string   `json:"dir"`
		Ext []string `json:"ext"`
	} `json:"content"`
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
