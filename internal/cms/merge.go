package cms

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/moonwalker/moonbase/pkg/content"
)

func MergeLocalisedContent(rc []*github.RepositoryContent, cs content.Schema) (*content.MergedContentData, error) {
	result := &content.MergedContentData{}
	fields := make(map[string]map[string]interface{})
	df := make(map[string]interface{})
	lf := make(map[string]bool)

	//Build Localized field bool from schema
	for _, csf := range cs.Fields {
		if csf.Localized {
			lf[csf.ID] = csf.Localized
		}
	}

	// Get default locale content
	for _, c := range rc {
		if *c.Name != content.JsonSchemaName {
			n, l := GetNameLocaleFromFilename(*c.Name)
			if content.DefaultLocale == l {
				result.Name = n
				result.Type = *c.Type
				cnt, err := c.GetContent()
				if err != nil {
					return nil, fmt.Errorf("error getting repo content: %s", err)
				}
				err = json.Unmarshal([]byte(cnt), &df)
				if err != nil {
					return nil, fmt.Errorf("error parsing localised content data: %s", err)
				}

				for k, v := range df {
					if k == "fields" {
						im, ok := v.(map[string]interface{})
						if !ok {
							return nil, fmt.Errorf("error parsing inner map[string]interface{}")
						}
						for ik, iv := range im {
							if fields[ik] == nil {
								fields[ik] = make(map[string]interface{})
							}
							if fields[ik][content.DefaultLocale] == nil {
								fields[ik][content.DefaultLocale] = make(map[string]interface{})
							}
							fields[ik][content.DefaultLocale] = iv
						}
					}
					// TODO: Do we need to include non Fields property data (e.g. createdAt, createdBy...)?
				}
				result.Fields = fields
			}
		}

	}

	// Loop through locales, skipping default and get content
	for _, c := range rc {
		if *c.Name != content.JsonSchemaName {
			lcd := &content.ContentData{}
			_, lcl := GetNameLocaleFromFilename(*c.Name)

			if lcl != content.DefaultLocale {
				cnt, err := c.GetContent()
				if err != nil {
					return nil, fmt.Errorf("error getting repo content: %s", err)
				}

				err = json.Unmarshal([]byte(cnt), &lcd)
				if err != nil {
					return nil, fmt.Errorf("error parsing localised content data: %s", err)
				}
			}

			for k, f := range lcd.Fields {
				if fields[k] == nil {
					fields[k] = make(map[string]interface{})
				}
				// TODO: Check if field is equal to default before adding?
				if _, ok := lf[k]; ok {
					fields[k][lcl] = f
				}
			}
		}
	}
	return result, nil
}

func GetNameLocaleFromFilename(fn string) (string, string) {
	ui := strings.LastIndex(fn, "_")
	di := strings.LastIndex(fn, ".")
	name := fn[:ui]
	locale := fn[ui+1 : di]

	return name, locale
}
