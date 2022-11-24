package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/gosimple/slug"

	"github.com/moonwalker/moonbase/internal/cms"
	"github.com/moonwalker/moonbase/internal/gh"
)

const (
	cmsConfigPath = "moonbase.yaml"
)

type collectionPayload struct {
	Name string `json:"name"`
}

type entryPayload struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

// config

func getConfig(ctx context.Context, accessToken string, owner string, repo string, ref string) *cms.Config {
	data, _ := gh.GetBlob(ctx, accessToken, owner, repo, ref, cmsConfigPath)
	return cms.ParseConfig(cmsConfigPath, data)
}

// dashboard

func getDash(w http.ResponseWriter, r *http.Request) {
	//
	w.WriteHeader(http.StatusOK)
}

// collections

// @Summary		Get collections
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Success		200	{object}	[]treeItem
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections	[get]
// @Security	bearerToken
func getCollections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	repoContents, err := gh.GetTree(ctx, accessToken, owner, repo, ref, cmsConfig.Content.Dir)
	if err != nil {
		errClientFailGetTree().Log(err).Json(w)
		return
	}

	treeItems := make([]*treeItem, 0)
	for _, rc := range repoContents {
		if *rc.Type == "dir" {
			treeItems = append(treeItems, &treeItem{
				Name: rc.Name,
				// Path: rc.Path,
				Type: rc.Type,
				SHA:  rc.SHA,
			})
		}
	}

	jsonResponse(w, http.StatusOK, treeItems)
}

// @Summary		Create or Update collection
// @Tags		cms
// @Accept		json
// @Param		owner			path	string		true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string		true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string		true	"git ref (branch, tag, sha)"
// @Param		payload	body	collectionPayload	true	"collection payload"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections	[post]
// @Security	bearerToken
func postCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	collection := &collectionPayload{}
	err := json.NewDecoder(r.Body).Decode(collection)
	if err != nil {
		errFailedDecReqBody().Log(err).Json(w)
		return
	}

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	collectionName := slug.Make(collection.Name)
	path := filepath.Join(cmsConfig.Content.Dir, collectionName, ".gitkeep")
	commitMessage := fmt.Sprintf("feat(content): create %s", collectionName)
	emptyContent := ""

	err = gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &emptyContent, commitMessage)
	if err != nil {
		errClientFailCommitBlob().Log(err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Delete collection
// @Tags		cms
// @Accept		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		collection		path	string	true	"collection"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection} [delete]
// @Security	bearerToken
func delCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	collectionName := slug.Make(collection)
	path := filepath.Join(cmsConfig.Content.Dir, collectionName, ".gitkeep")
	commitMessage := fmt.Sprintf("feat(content): delete %s", collectionName)

	err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, nil, commitMessage)
	if err != nil {
		errClientFailDeleteBlob().Log(err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// entries

// @Summary		Get entries
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		collection		path	string	true	"collection"
// @Success		200	{object}	[]treeItem
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection}	[get]
// @Security	bearerToken
func getEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.Content.Dir, collection)

	repoContents, err := gh.GetTree(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errClientFailGetTree().Log(err).Json(w)
		return
	}

	treeItems := make([]*treeItem, 0)
	for _, rc := range repoContents {
		if *rc.Type == "file" {
			treeItems = append(treeItems, &treeItem{
				Name: rc.Name,
				// Path: rc.Path,
				Type: rc.Type,
				SHA:  rc.SHA,
			})
		}
	}

	jsonResponse(w, http.StatusOK, treeItems)
}

// @Summary		Create or Update entry
// @Tags		cms
// @Accept		json
// @Param		owner			path	string			true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string			true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string			true	"git ref (branch, tag, sha)"
// @Param		collection		path	string			true	"collection"
// @Param		payload			body	entryPayload	true	"entry payload"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection}	[post]
// @Security	bearerToken
func postEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")

	entry := &entryPayload{}
	err := json.NewDecoder(r.Body).Decode(entry)
	if err != nil {
		errFailedDecReqBody().Log(err).Json(w)
		return
	}

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	entryName := slug.Make(entry.Name)
	path := filepath.Join(cmsConfig.Content.Dir, collection, entryName)
	commitMessage := fmt.Sprintf("feat(%s): create/update %s", collection, entryName)

	err = gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &entry.Contents, commitMessage)
	if err != nil {
		errClientFailCommitBlob().Log(err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Get entry
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		collection		path	string	true	"collection"
// @Param		entry			path	string	true	"entry"
// @Success		200	{object}	entryPayload
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/{collection}/{entry}	[get]
// @Security	bearerToken
func getEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")
	entry := chi.URLParam(r, "entry")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.Content.Dir, collection, entry)

	blob, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errClientFailGetBlob().Log(err).Json(w)
		return
	}

	data := &blobEntry{blob}
	jsonResponse(w, http.StatusOK, data)
}

// @Summary		Delete entry
// @Tags		cms
// @Accept		json
// @Param		owner			path	string		true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string		true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string		true	"git ref (branch, tag, sha)"
// @Param		collection		path	string		true	"collection"
// @Param		entry			path	string		true	"entry"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}	[delete]
// @Security	bearerToken
func delEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")
	entry := chi.URLParam(r, "entry")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.Content.Dir, collection, entry)
	commitMessage := fmt.Sprintf("delete(%s): %s", collection, entry)

	err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, nil, commitMessage)
	if err != nil {
		errClientFailDeleteBlob().Log(err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
