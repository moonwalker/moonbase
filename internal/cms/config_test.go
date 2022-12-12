package cms

import (
	"errors"
	"os"
	"testing"
)

const (
	yamlPath = "testdata/moonbase.yaml"
	jsonPath = "testdata/moonbase.json"
)

func TestConfigParseYAML(t *testing.T) {
	testParse(t, yamlPath, "test")
}

func TestConfigParseJSON(t *testing.T) {
	testParse(t, jsonPath, "")
}

func testParse(t *testing.T, path string, dir string) {
	data, _ := os.ReadFile(path)
	config := ParseConfig(data)
	if config.WorkDir != dir {
		t.Error(errors.New("working dir mismatch"))
	}
}
