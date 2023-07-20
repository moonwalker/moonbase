package cms

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/moonwalker/moonbase/pkg/content"
)

func MergeLocalisedContent(rc []*github.RepositoryContent) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	fmt.Println("RC count:", len(rc))
	for _, c := range rc {
		if *c.Name != content.JsonSchemaName {
			cd := &content.ContentData{}
			err := json.Unmarshal([]byte(*c.Content), cd)
			if err != nil {
				return nil, errors.New("error parsing localised content data")
			}
			fmt.Println("ContentData:", cd)

			fmt.Println("CD Fields count:", len(cd.Fields))
			for _, f := range cd.Fields {
				fieldJSON, err := json.Marshal(f)
				if err != nil {
					return nil, fmt.Errorf("error marshalling field data: %s", err)
				}
				fmt.Println("CD Fields:", f)
				fmt.Println("FieldJSON:", fieldJSON)
				tempField := &content.Field{}
				err = json.Unmarshal(fieldJSON, &tempField)
				if err != nil {
					return nil, fmt.Errorf("error unmarshalling field data: %s", err)
				}
				// if contains(lfs, tempField.ID) {
				// 	fmt.Println("Temp field name matched", tempField.ID)
				// }
			}
			// if len(result) == 0 {
			// 	blobType := filepath.Ext(*c.Name)
			// 	result, err = ParseBlob(blobType, blob)
			// 	if err != nil {
			// 		return nil, errors.New("error parsing content blob")
			// 	}

			// 	fmt.Println("Result:", result)
			fmt.Println("Locale:", getLocaleFromFilename(*c.Name))
			// }
		}
	}
	// fmt.Println("Result:", result)

	return result, nil
}

func getLocaleFromFilename(fn string) string {
	ui := strings.LastIndex(fn, "_")
	di := strings.LastIndex(fn, ".")
	locale := fn[ui+1 : di]

	return locale
}

func contains(s []string, v string) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}
