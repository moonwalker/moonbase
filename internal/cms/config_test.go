package cms

import (
	"errors"
	"os"
	"testing"
)

const (
	contentdir = "content"
	jsonPath   = "testdata/moonbase.json"
	yamlPath   = "testdata/moonbase.yaml"
)

func TestConfigParseJSON(t *testing.T) {
	testParse(t, jsonPath)
}

func TestConfigParseYAML(t *testing.T) {
	testParse(t, yamlPath)
}

func testParse(t *testing.T, path string) {
	data, _ := os.ReadFile(path)
	config := ParseConfig(path, data)
	if config.ContentDir != contentdir {
		t.Error(errors.New("content dir mismatch"))
	}
}
