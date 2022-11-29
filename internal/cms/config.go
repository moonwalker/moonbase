package cms

import (
	"encoding/json"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ContentDir string      `json:"contents" yaml:"contents"`
	Components compsConfig `json:"components" yaml:"components"`
}

type compsConfig struct {
	Entry        string   `json:"entry" yaml:"entry"`
	Dependencies []string `json:"dependencies" yaml:"dependencies"`
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

// entry can be dir or file, this always returns the dir
func (c *compsConfig) EntryDir() string {
	if len(filepath.Ext(c.Entry)) > 0 {
		return filepath.Dir(c.Entry)
	}
	return c.Entry
}
