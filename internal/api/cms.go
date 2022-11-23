package api

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/internal/cms"
	"github.com/moonwalker/moonbase/internal/gh"
)

const (
	cmsConfigPath = "moonbase.yaml"
)

type collectionPayload struct {
	Name  string `json:"name"`
	User  string `json:"user"`
	Email string `json:"email"`
}

type documentPayload struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Email    string `json:"email"`
	Contents []byte `json:"contents"`
}

// config

func getConfig(ctx context.Context, accessToken string, owner string, repo string, ref string) *cms.Config {
	data, _ := gh.GetBlob(ctx, accessToken, owner, repo, ref, cmsConfigPath)
	return cms.ParseConfig(cmsConfigPath, data)
}

// collection

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
				Path: rc.Path,
				Type: rc.Type,
				SHA:  rc.SHA,
			})
		}
	}

	jsonResponse(w, http.StatusOK, treeItems)
}

func newCollection(w http.ResponseWriter, r *http.Request) {
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
	path := filepath.Join(cmsConfig.Content.Dir, collection.Name)

	err = gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, collection.User, collection.Email, "")
	if err != nil {
		errClientFailCommitBlob().Log(err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// documents

func getDocuments(w http.ResponseWriter, r *http.Request) {
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
		if *rc.Type == "file" {
			treeItems = append(treeItems, &treeItem{
				Name: rc.Name,
				Path: rc.Path,
				Type: rc.Type,
				SHA:  rc.SHA,
			})
		}
	}

	jsonResponse(w, http.StatusOK, treeItems)
}

func newDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	document := &documentPayload{}
	err := json.NewDecoder(r.Body).Decode(document)
	if err != nil {
		errFailedDecReqBody().Log(err).Json(w)
		return
	}

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.Content.Dir, document.Name)

	err = gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, document.User, document.Email, string(document.Contents))
	if err != nil {
		errClientFailCommitBlob().Log(err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	document := chi.URLParam(r, "document")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.Content.Dir, document)

	blob, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errClientFailGetBlob().Log(err).Json(w)
		return
	}

	data := &blobEntry{blob}
	jsonResponse(w, http.StatusOK, data)
}
