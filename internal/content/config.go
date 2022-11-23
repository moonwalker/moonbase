package content

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ContentConfig struct {
	Content struct {
		Dir string   `json:"dir"`
		Ext []string `json:"ext"`
	} `json:"content"`
}

func ParseConfig(path string, data []byte) (*ContentConfig, bool) {
	config, err := unmarshalConfig(path, data)
	if err != nil {
		// log error
		return nil, false
	}
	return config, true
}

func unmarshalConfig(path string, data []byte) (*ContentConfig, error) {
	switch filepath.Ext(path) {
	case ".yaml":
		cfg := &ContentConfig{}
		err := yaml.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	case ".json":
		cfg := &ContentConfig{}
		err := json.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return nil, fmt.Errorf("unsupported config file type")
}
