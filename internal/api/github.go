package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/internal/gh"
)

type repositoryList struct {
	LastPage int               `json:"lastPage"`
	Items    []*repositoryItem `json:"items"`
}

type repositoryItem struct {
	Name  *string `json:"name"`
	Owner *string `json:"owner"`
}

type branchItem struct {
	Name *string `json:"name"`
	SHA  *string `json:"sha"`
}

type treeItem struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
	SHA  *string `json:"sha"`
}

const (
	configYamlPath = "moonbase.yaml"
)

// @Summary		List repositories
// @Tags		repositories
// @Accept		json
// @Produce		json
// @Param		page			query	string	false	"page of results to retrieve (default: `1`)"
// @Param		perpage			query	string	false	"number of results to include per page (default: `30`)"
// @Param		sort			query	string	false	"how to sort the repository list, can be one of `created`, `updated`, `pushed`, `full_name` (default: `full_name`)"
// @Param		direction		query	string	false	"direction in which to sort repositories, can be one of `asc` or `desc` (default when using `full_name`: `asc`; otherwise: `desc`)"
// @Success		200	{object}	repositoryList
// @Failure		500	{object}	errorData
// @Router		/list [get]
// @Security	bearerToken
func getRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)

	page, _ := strconv.Atoi(r.FormValue("page"))
	perpage, _ := strconv.Atoi(r.FormValue("perpage"))
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")

	grs, lastPage, err := gh.ListRepositories(ctx, accessToken, page, perpage, sort, direction)
	if err != nil {
		errClientFailGetRepositories().Log(err).Json(w)
		return
	}

	repoItems := make([]*repositoryItem, 0)
	for _, gr := range grs {
		repoItems = append(repoItems, &repositoryItem{
			Name:  gr.Name,
			Owner: gr.Owner.Login,
		})
	}

	repoList := &repositoryList{LastPage: lastPage, Items: repoItems}
	jsonResponse(w, http.StatusOK, repoList)
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
