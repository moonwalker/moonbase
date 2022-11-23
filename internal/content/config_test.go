package content

import (
	"os"
	"testing"
)

const (
	contentDir = "content"
	jsonPath   = "moonbase.json"
	yamlPath   = "moonbase.yaml"
)

func TestConfigParseJSON(t *testing.T) {
	testParse(t, jsonPath)
}

func TestConfigParseYAML(t *testing.T) {
	testParse(t, yamlPath)
}

func testParse(t *testing.T, path string) {
	data, _ := os.ReadFile(path)
	config, ok := ParseConfig(path, data)
	if !ok {
		t.Fail()
	}
	if config.Content.Dir != contentDir {
		t.Fail()
	}
}
