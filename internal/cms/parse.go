package cms

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	markdownTpl = "---\n%s---\n%s"
	dateFormat  = "2006-01-02"
)

func ParseBlob(contentType string, content string) (map[string]interface{}, error) {
	switch contentType {
	case ".md", ".mdx":
		return parseMarkdown(content)
	case ".json":
		return ParseJSON(content)
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

func ParseJSON(content string) (map[string]interface{}, error) {
	mc := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &mc)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

func JsonToMarkdown(json string) (string, error) {
	pj, err := ParseJSON(json)
	if err != nil {
		return "", err
	}

	body := pj["body"]
	delete(pj, "body")

	for k, v := range pj {
		sv, ok := v.(string)
		if ok {
			t := isDateValue(sv)
			if !t.IsZero() {
				pj[k] = t
			}
		}
	}

	ym, err := yaml.Marshal(pj)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(markdownTpl, string(ym), body), nil
}

func isDateValue(stringDate string) time.Time {
	t, _ := time.Parse(dateFormat, stringDate)
	return t
}
