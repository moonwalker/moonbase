package cms

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/moonwalker/moonbase/pkg/content"
	gh "github.com/moonwalker/moonbase/pkg/github"
)

type localePayloads struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

func GetNameLocaleFromPath(path string) (string, string) {
	dirs := strings.Split(filepath.Dir(path), "/")
	name := dirs[len(dirs)-1]
	locale := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

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
			n, l := GetNameLocaleFromPath(*c.Path)

			cnt, err := c.GetContent()
			if err != nil {
				return nil, fmt.Errorf("error getting repo content: %s", err)
			}

			dcd := &content.ContentData{}
			err = json.Unmarshal([]byte(cnt), dcd)
			if err != nil {
				return nil, fmt.Errorf("error parsing localised content data: %s", err)
			}

			if l == content.DefaultLocale {
				result.Name = n
				result.Type = *c.Type
				result.ID = dcd.ID
				result.Version = dcd.Version

				if dcd.CreatedAt != "" {
					ct, _ := time.Parse(time.RFC3339Nano, dcd.CreatedAt)
					result.CreatedAt = &ct
				}
				result.CreatedBy = dcd.CreatedBy
				if dcd.UpdatedAt != "" {
					ut, _ := time.Parse(time.RFC3339Nano, dcd.UpdatedAt)
					result.UpdatedAt = &ut
				}
				result.UpdatedBy = dcd.UpdatedBy
			}

			for k, v := range dcd.Fields {
				if result.Fields[k] == nil {
					result.Fields[k] = make(map[string]interface{})
				}

				if l == content.DefaultLocale {
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

func SeparateLocalisedContent(user string, mcd content.MergedContentData, locales []string, workDir, collection string) ([]gh.BlobEntry, error) {
	var res []gh.BlobEntry
	now := time.Now().UTC().Format(time.RFC3339Nano)

	for _, l := range locales {
		fields := make(map[string]interface{})
		for key, value := range mcd.Fields {
			if fields[key] == nil {
				fields[key] = make(map[string]interface{})
			}
			if (value[l] != nil && value[l] != "") || l == content.DefaultLocale {
				fields[key] = value[l]
			} else {
				fields[key] = value[content.DefaultLocale]
			}
		}

		s, err := json.Marshal(content.ContentData{
			ID:        mcd.ID,
			CreatedAt: mcd.CreatedAt.Format(time.RFC3339Nano),
			CreatedBy: mcd.CreatedBy,
			UpdatedAt: now,
			UpdatedBy: user,
			Fields:    fields,
		})
		if err != nil {
			return nil, fmt.Errorf("error marshalling content data:%s", err)
		}

		content := string(s)
		res = append(res, gh.BlobEntry{
			Path:    filepath.Join(workDir, collection, mcd.Name, fmt.Sprintf("%s.json", l)),
			Content: &content,
		})
	}

	return res, nil
}
