package api

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/moonwalker/moonbase/internal/cms"
	gh "github.com/moonwalker/moonbase/pkg/github"
)

const localesConfig = "locales.json"

func getLocales(ctx context.Context, accessToken, owner, repo, ref string) ([]string, int, error) {
	path := filepath.Join(cms.SettingsFolder, localesConfig)

	fc, resp, err := gh.GetFileContent(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	blob, err := fc.GetContent()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	locales := make([]string, 0)
	err = json.Unmarshal([]byte(blob), &locales)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return locales, http.StatusOK, nil
}

// func getLocales(ctx context.Context, accessToken, owner, repo, branch, path string) ([]string, int, error) {
// 	rcs, resp, err := gh.GetAllLocaleContents(ctx, accessToken, owner, repo, branch, path)
// 	if err != nil {
// 		return nil, resp.StatusCode, err
// 	}

// 	res := make([]string, 0)
// 	for _, rc := range rcs {
// 		if *rc.Type == "file" && *rc.Name != content.JsonSchemaName {
// 			fn, l := cms.GetNameLocaleFromPath(*rc.Path)
// 			if fn != "" {
// 				res = append(res, l)
// 			}
// 		}
// 	}
// 	return res, 0, nil
// }
