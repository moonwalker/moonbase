package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"

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
	Path *string `json:"path,omitempty"`
	Type *string `json:"type"`
	SHA  *string `json:"sha"`
}

type blobEntry struct {
	Contents []byte `json:"contents"`
}

type commitPayload struct {
	Contents      []byte `json:"contents"`
	CommitMessage string `json:"commitMessage"`
}

type componentItem struct {
	Path    *string `json:"path"`
	Content []byte  `json:"content"`
}

// @Summary		Get repositories
// @Tags		repos
// @Accept		json
// @Produce		json
// @Param		page			query	string	false	"page of results to retrieve (default: `1`)"
// @Param		per_page		query	string	false	"number of results to include per page (default: `30`)"
// @Param		sort			query	string	false	"how to sort the repository list, can be one of `created`, `updated`, `pushed`, `full_name` (default: `full_name`)"
// @Param		direction		query	string	false	"direction in which to sort repositories, can be one of `asc` or `desc` (default when using `full_name`: `asc`; otherwise: `desc`)"
// @Success		200	{object}	repositoryList
// @Failure		500	{object}	errorData
// @Router		/repos [get]
// @Security	bearerToken
func getRepos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	page, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("per_page"))
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")

	grs, resp, err := gh.ListRepositories(ctx, accessToken, page, perPage, sort, direction)
	if err != nil {
		errReposGet().Status(resp.StatusCode).Log(r, err).Json(w)
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

	repoList := &repositoryList{LastPage: resp.LastPage, Items: repoItems}
	jsonResponse(w, http.StatusOK, repoList)
}

// @Summary		Get branhces
// @Tags		repos
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Success		200	{object}	branchList
// @Failure		500	{object}	errorData
// @Router		/repos/{owner}/{repo}/branches [get]
// @Security	bearerToken
func getBranches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	branches, resp, err := gh.ListBranches(ctx, accessToken, owner, repo)
	if err != nil {
		errReposGetBranches().Status(resp.StatusCode).Log(r, err).Json(w)
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

// @Summary		Get tree
// @Tags		repos
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		path			path	string	true	"tree path"
// @Success		200	{object}	[]byte
// @Failure		500	{object}	errorData
// @Router		/repos/{owner}/{repo}/tree/{ref}/{path} [get]
// @Security	bearerToken
func getTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	path := chi.URLParam(r, "*")

	repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errReposGetTree().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	treeItems := make([]*treeItem, 0)
	for _, rc := range repoContents {
		treeItems = append(treeItems, &treeItem{
			Name: rc.Name,
			Path: rc.Path,
			Type: rc.Type,
			SHA:  rc.SHA,
		})
	}

	jsonResponse(w, http.StatusOK, treeItems)
}

// @Summary		Get blob
// @Tags		repos
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		path			path	string	true	"contents path"
// @Success		200	{object}	blobEntry
// @Failure		500	{object}	errorData
// @Router		/repos/{owner}/{repo}/blob/{ref}/{path} [get]
// @Security	bearerToken
func getBlob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	path := chi.URLParam(r, "*")

	blob, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		e := errReposGetBlob()
		if resp != nil {
			e.Status(resp.StatusCode)
		}
		e.Log(r, err).Json(w)
		return
	}

	data := &blobEntry{blob}
	jsonResponse(w, http.StatusOK, data)
}

// @Summary		Post blob
// @Tags		repos
// @Accept		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		path			path	string	true	"contents path"
// @Param		payload			body	commitPayload	true	"commit payload"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/repos/{owner}/{repo}/blob/{ref}/{path} [post]
// @Security	bearerToken
func postBlob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	path := chi.URLParam(r, "*")

	data := &commitPayload{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		errJsonDecode().Log(r, err).Json(w)
		return
	}

	contents := string(data.Contents)
	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &contents, string(data.CommitMessage))
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Delete blob
// @Tags		repos
// @Accept		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		path			path	string	true	"contents path"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/repos/{owner}/{repo}/blob/{ref}/{path} [delete]
// @Security	bearerToken
func delBlob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	path := chi.URLParam(r, "*")

	deleteMessage := fmt.Sprintf("delete %s", filepath.Base(path))
	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, nil, deleteMessage)
	if err != nil {
		errReposDeleteBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
