package cms

import (
	"os"
	"testing"
)

var payload = map[string]string{
	"": "",
}

func TestBundleComponents(t *testing.T) {
	data, _ := os.ReadFile(yamlPath)
	config := ParseConfig(yamlPath, data)
	res, err := BundleComponents(config.Components.Entry, config.Components.Dependencies, false, false)
	if err != nil {
		t.Error(err)
	}
	println(res)
}
