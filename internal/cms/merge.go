package cms

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/moonwalker/moonbase/pkg/content"
)

func GetNameLocaleFromFilename(fn string) (string, string) {
	ui := strings.LastIndex(fn, "_")
	di := strings.LastIndex(fn, ".")
	name := fn[:ui]
	locale := fn[ui+1 : di]

	return name, locale
}

func MergeLocalisedContent(rc []*github.RepositoryContent, cs content.Schema) (*content.MergedContentData, error) {
	result := &content.MergedContentData{}
	result.Fields = make(map[string]map[string]interface{})

	localizedFields := make(map[string]bool)

	//Build Localized field bool from schema
	for _, csf := range cs.Fields {
		if csf.Localized {
			localizedFields[csf.ID] = csf.Localized
		}
	}

	// Get default locale content
	for _, c := range rc {
		if *c.Name != content.JsonSchemaName {
			n, l := GetNameLocaleFromFilename(*c.Name)

			cnt, err := c.GetContent()
			if err != nil {
				return nil, fmt.Errorf("error getting repo content: %s", err)
			}
			dcd := &content.ContentData{}
			err = json.Unmarshal([]byte(cnt), dcd)
			if err != nil {
				return nil, fmt.Errorf("error parsing localised content data: %s", err)
			}

			for k, v := range dcd.Fields {
				if result.Fields[k] == nil {
					result.Fields[k] = make(map[string]interface{})
				}

				if l == content.DefaultLocale {
					result.Name = n
					result.Type = *c.Type

					result.ID = dcd.ID
					if dcd.CreatedAt != "" {
						ct, _ := time.Parse(time.RFC3339, dcd.CreatedAt)
						result.CreatedAt = &ct
					}
					result.CreatedBy = dcd.CreatedBy
					if dcd.UpdatedAt != "" {
						ut, _ := time.Parse(time.RFC3339, dcd.UpdatedAt)
						result.UpdatedAt = &ut
					}
					result.UpdatedBy = dcd.UpdatedBy
					result.Version = dcd.Version

					result.Fields[k][l] = v
					if localizedFields[k] && len(result.Fields[k]) > 0 {
						// clear existing locale values if they are the same as the default if any
						for el, ev := range result.Fields[k] {
							if el != content.DefaultLocale && ev == v {
								result.Fields[k][el] = nil
							}
						}
					}
				} else if localizedFields[k] {
					// check if field is localized and locale value is not the same as the default
					if v != result.Fields[k][content.DefaultLocale] {
						result.Fields[k][l] = v
					} else {
						result.Fields[k][l] = nil
					}
				}
			}
		}
	}

	return result, nil
}
