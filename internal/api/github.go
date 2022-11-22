package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/internal/gh"
)

type branchItem struct {
	Name *string `json:"name"`
	SHA  *string `json:"sha"`
}

type repositoryItem struct {
	Name  *string `json:"name"`
	Owner *string `json:"owner"`
}

type treeItem struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
	SHA  *string `json:"sha"`
}

const (
	configYamlPath = "moonbase.yaml"
)

// @Summary	List repositories
// @Tags		github
// @Accept		json
// @Produce	json
// @Success	200	{array} listItem "ok"
// @Router		/list [get]
// @Security	bearerToken
func getRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)

	grs, err := gh.ListRepositories(ctx, accessToken, 100, 1)
	if err != nil {
		errClientFailGetRepositories().Log(err).Json(w)
		return
	}

	repos := make([]*repositoryItem, 0)
	for _, gr := range grs {
		repos = append(repos, &repositoryItem{
			Name:  gr.Name,
			Owner: gr.Owner.Login,
		})
	}

	jsonResponse(w, http.StatusOK, repos)
}

func getBranches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	branches, err := gh.ListBranches(ctx, accessToken, owner, repo)
	if err != nil {
		errClientFailGetBranches().Log(err).Json(w)
		return
	}

	bs := make([]*branchItem, 0)
	for _, br := range branches {
		bi := branchItem{
			Name: br.Name,
		}
		if br.Commit != nil {
			bi.SHA = br.Commit.SHA
		}
		bs = append(bs, &bi)
	}

	jsonResponse(w, http.StatusOK, bs)
}
