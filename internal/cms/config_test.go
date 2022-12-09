package cms

import (
	"errors"
	"os"
	"testing"
)

const (
	workdir  = "content"
	jsonPath = "testdata/moonbase.json"
	yamlPath = "testdata/moonbase.yaml"
)

func TestConfigParseJSON(t *testing.T) {
	testParse(t, jsonPath)
}

func TestConfigParseYAML(t *testing.T) {
	testParse(t, yamlPath)
}

func testParse(t *testing.T, path string) {
	data, _ := os.ReadFile(path)
	config := ParseConfig(data)
	if config.WorkDir != workdir {
		t.Error(errors.New("working dir mismatch"))
	}
}
