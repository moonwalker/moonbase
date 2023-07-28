package cms

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/moonwalker/moonbase/pkg/content"
)

func MergeLocalisedContent(rc []*github.RepositoryContent) (map[string]interface{}, error) {
	mcd := &content.MergedContentData{}

	for _, c := range rc {
		if *c.Name != content.JsonSchemaName {
			lc := &content.LocalisedContentData{
				Locale: getLocaleFromFilename(*c.Name),
			}
			err := json.Unmarshal([]byte(*c.Content), &lc.ContentData)
			if err != nil {
				return nil, errors.New("error parsing localised content data")
			}

			mcd.LocalisedContent = append(mcd.LocalisedContent, lc)
		}
	}
	result := map[string]interface{}{
		"Slug":             mcd.Slug,
		"LocalisedContent": mcd.LocalisedContent,
	}
	return result, nil
}

func getLocaleFromFilename(fn string) string {
	ui := strings.LastIndex(fn, "_")
	di := strings.LastIndex(fn, ".")
	locale := fn[ui+1 : di]

	return locale
}
