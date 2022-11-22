package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"encoding/json"

	"github.com/moonwalker/moonbase/internal/gh"
)

type repositoryList struct {
	LastPage int               `json:"lastPage"`
	Items    []*repositoryItem `json:"items"`
}

type repositoryItem struct {
	Name          *string `json:"name"`
	Owner         *string `json:"owner"`
	DefaultBranch *string `json:"defaultBranch"`
}

type branchList struct {
	Items []*branchItem `json:"items"`
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
// @Param		per_page		query	string	false	"number of results to include per page (default: `30`)"
// @Param		sort			query	string	false	"how to sort the repository list, can be one of `created`, `updated`, `pushed`, `full_name` (default: `full_name`)"
// @Param		direction		query	string	false	"direction in which to sort repositories, can be one of `asc` or `desc` (default when using `full_name`: `asc`; otherwise: `desc`)"
// @Success		200	{object}	repositoryList
// @Failure		500	{object}	errorData
// @Router		/user/repos [get]
// @Security	bearerToken
func getRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)

	page, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("per_page"))
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")

	grs, lastPage, err := gh.ListRepositories(ctx, accessToken, page, perPage, sort, direction)
	if err != nil {
		errClientFailGetRepositories().Log(err).Json(w)
		return
	}

	repoItems := make([]*repositoryItem, 0)
	for _, gr := range grs {
		repoItems = append(repoItems, &repositoryItem{
			Name:          gr.Name,
			Owner:         gr.Owner.Login,
			DefaultBranch: gr.DefaultBranch,
		})
	}

	repoList := &repositoryList{LastPage: lastPage, Items: repoItems}
	jsonResponse(w, http.StatusOK, repoList)
}

// @Summary		List branhces
// @Tags		branches
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Success		200	{object}	branchList
// @Failure		500	{object}	errorData
// @Router		/{owner}/{repo}/branches [get]
// @Security	bearerToken
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

	branchItems := make([]*branchItem, 0)
	for _, br := range branches {
		branchItem := branchItem{
			Name: br.Name,
		}
		if br.Commit != nil {
			branchItem.SHA = br.Commit.SHA
		}
		branchItems = append(branchItems, &branchItem)
	}

	branchList := &branchList{Items: branchItems}
	jsonResponse(w, http.StatusOK, branchList)
}

func getTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	branch := chi.URLParam(r, "branch")
	path := chi.URLParam(r, "*")

	tree, err := gh.GetTree(ctx, accessToken, owner, repo, branch, path)
	if err != nil {
		errClientFailGetTree().Log(err).Json(w)
		return
	}

	treeItems := make([]*treeItem, 0)
	for _, te := range tree.Entries {
		treeItems = append(treeItems, &treeItem{
			Name: te.Path,
			Type: te.Type,
			SHA:  te.SHA,
		})
	}

	json.NewEncoder(w).Encode(treeItems)
}

// @Summary		Get blob
// @Tags		contents
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		path			path	string	true	"contents path"
// @Success		200	{object}	[]byte
// @Failure		500	{object}	errorData
// @Router		/repos/{owner}/{repo}/blob/{ref}/{path} [get]
// @Security	bearerToken
func getBlob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	path := chi.URLParam(r, "*")

	blob, err := gh.GetBlobByPath(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errClientFailGetBlob().Log(err).Json(w)
		return
	}

	data := struct {
		Contents []byte `json:"contents"`
	}{
		blob,
	}
	json.NewEncoder(w).Encode(data)
}
