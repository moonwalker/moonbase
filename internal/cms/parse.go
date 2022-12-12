package cms

import (
	"encoding/json"
	"errors"
	"regexp"

	"gopkg.in/yaml.v2"
)

func ParseBlob(contentType string, content string) (map[string]interface{}, error) {
	switch contentType {
	case ".md", ".mdx":
		return parseMarkdown(content)
	case ".json":
		return parseJSON(content)
	default:
		return nil, errors.New("not supported file format")
	}
}

func parseMarkdown(content string) (map[string]interface{}, error) {
	mc := make(map[string]interface{})
	r := regexp.MustCompile(`(^---(?:\r?\n|\r)([\s\S]*)---(?:\r?\n|\r)?)?([\s\S]*)`)
	matches := r.FindStringSubmatch(content)

	if len(matches) > 3 {
		err := yaml.Unmarshal([]byte(matches[2]), &mc)
		if err != nil {
			return nil, err
		}
		mc["body"] = matches[3]
	}

	return mc, nil
}

func parseJSON(content string) (map[string]interface{}, error) {
	mc := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &mc)
	if err != nil {
		return nil, err
	}
	return mc, nil
}
