package api

import (
	"context"

	"github.com/moonwalker/moonbase/internal/cms"
	"github.com/moonwalker/moonbase/pkg/content"
	gh "github.com/moonwalker/moonbase/pkg/github"
)

func getLocales(ctx context.Context, accessToken, owner, repo, branch, path string) ([]string, int, error) {
	rcs, resp, err := gh.GetAllLocaleContents(ctx, accessToken, owner, repo, branch, path)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	res := make([]string, 0)
	for _, rc := range rcs {
		if *rc.Type == "file" && *rc.Name != content.JsonSchemaName {
			fn, l := cms.GetNameLocaleFromPath(*rc.Path)
			if fn != "" {
				res = append(res, l)
			}
		}
	}
	return res, 0, nil
}
