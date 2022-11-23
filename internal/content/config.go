package content

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type CmsConfig struct {
	Content struct {
		Dir string `json:"dir"`
	} `json:"content"`
}

func ParseConfig(path string, data []byte) (*CmsConfig, bool) {
	config, err := unmarshalConfig(path, data)
	if err != nil {
		// log error
		return nil, false
	}
	return config, true
}

func unmarshalConfig(path string, data []byte) (*CmsConfig, error) {
	switch filepath.Ext(path) {
	case ".yaml":
		cfg := &CmsConfig{}
		err := yaml.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	case ".json":
		cfg := &CmsConfig{}
		err := json.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return nil, fmt.Errorf("unsupported config file type")
}