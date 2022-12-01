package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gosimple/slug"

	"github.com/moonwalker/moonbase/internal/cache"
	"github.com/moonwalker/moonbase/internal/cms"
	"github.com/moonwalker/moonbase/internal/gh"
)

type collectionPayload struct {
	Name string `json:"name"`
}

type entryPayload struct {
	Name       string `json:"name"`
	Contents   string `json:"contents"`
	SaveSchema bool   `json:"save_schema"`
}

type commitEntry struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

// config

func getConfig(ctx context.Context, accessToken string, owner string, repo string, ref string) *cms.Config {
	data, _, _ := gh.GetBlob(ctx, accessToken, owner, repo, ref, cms.ConfigPath)
	return cms.ParseConfig(data)
}

// schema

func getSchema(ctx context.Context, accessToken string, owner string, repo string, ref string, collection string, workdir string) *cms.Schema {
	p := filepath.Join(workdir, collection, cms.JsonSchemaName)
	data, _, _ := gh.GetBlob(ctx, accessToken, owner, repo, ref, p)
	return cms.NewSchema(data)
}

// info

// @Summary		Get info
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Success		200	{object}	[]commitEntry
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}	[get]
// @Security	bearerToken
func getInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	rc, resp, err := gh.GetCommits(ctx, accessToken, owner, repo, ref)
	if err != nil {
		errCmsGetCommits().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	changes := make([]*commitEntry, 0)
	for _, c := range rc {
		changes = append(changes, &commitEntry{
			Author:  *c.Commit.Author.Name,
			Message: *c.Commit.Message,
			Date:    c.Commit.Author.Date.UTC().String(),
		})
	}

	jsonResponse(w, http.StatusOK, changes)
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

	repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, cmsConfig.ContentDir)
	if err != nil {
		errReposGetTree().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	treeItems := make([]*treeItem, 0)
	for _, rc := range repoContents {
		if *rc.Type == "dir" {
			treeItems = append(treeItems, &treeItem{
				Name: rc.Name,
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
		errJsonDecode().Log(r, err).Json(w)
		return
	}

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	collectionName := slug.Make(collection.Name)
	path := filepath.Join(cmsConfig.ContentDir, collectionName, ".gitkeep")
	commitMessage := fmt.Sprintf("feat(content): create %s", collectionName)
	emptyContent := ""

	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &emptyContent, commitMessage)
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
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
	path := filepath.Join(cmsConfig.ContentDir, collectionName)
	if path == cmsConfig.ContentDir {
		m := "missing collection name"
		errCmsDeleteFolder().Details(m).Log(r, errors.New(m)).Json(w)
		return
	}

	commitMessage := fmt.Sprintf("feat(content): delete %s", collectionName)
	resp, err := gh.DeleteFolder(ctx, accessToken, owner, repo, ref, path, commitMessage)
	if err != nil {
		errCmsDeleteFolder().Status(resp.StatusCode).Log(r, err).Json(w)
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
	path := filepath.Join(cmsConfig.ContentDir, collection)

	repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errReposGetTree().Status(resp.StatusCode).Log(r, err).Json(w)
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

// @Summary		Create entry
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
	createOrUpdateEntry(w, r)
}

// @Summary		Update entry
// @Tags		cms
// @Accept		json
// @Param		owner			path	string			true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string			true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string			true	"git ref (branch, tag, sha)"
// @Param		collection		path	string			true	"collection"
// @Param		entry			path	string			true	"entry"
// @Param		payload			body	entryPayload	true	"entry payload"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}	[put]
// @Security	bearerToken
func putEntry(w http.ResponseWriter, r *http.Request) {
	createOrUpdateEntry(w, r)
}

func createOrUpdateEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")
	entry := chi.URLParam(r, "entry")

	entryData := &entryPayload{}
	err := json.NewDecoder(r.Body).Decode(entryData)
	if err != nil {
		errJsonDecode().Log(r, err).Json(w)
		return
	}

	if len(entryData.Name) == 0 {
		entryData.Name = entry
	}
	if len(entryData.Name) == 0 {
		m := "missing entry name"
		errReposCommitBlob().Details(m).Log(r, errors.New(m)).Json(w)
	}

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	if !entryData.SaveSchema {
		schema := getSchema(ctx, accessToken, owner, repo, ref, collection, cmsConfig.ContentDir)
		if schema.Available() {
			err = schema.ValidateString(entryData.Contents)
			if err != nil {
				errCmsSchemaValidation().Log(r, err).Json(w)
				return
			}
		}
	}

	ext := filepath.Ext(entryData.Name)
	fn := strings.TrimSuffix(filepath.Base(entryData.Name), ext)
	entryName := slug.Make(fn) + ext

	path := filepath.Join(cmsConfig.ContentDir, collection, entryName)
	commitMessage := fmt.Sprintf("feat(%s): create/update %s", collection, entryName)

	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &entryData.Contents, commitMessage)
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	if entryData.SaveSchema {
		schema, err := cms.GenerateSchema(entryData.Contents)
		if err != nil {
			errCmsSchemaGeneration().Log(r, err).Json(w)
			return
		}
		schemaPath := filepath.Join(cmsConfig.ContentDir, collection, cms.JsonSchemaName)
		schemaCommitMessage := fmt.Sprintf("feat(%s): create/update %s", collection, cms.JsonSchemaName)
		resp, err = gh.CommitBlob(ctx, accessToken, owner, repo, ref, schemaPath, &schema, schemaCommitMessage)
		if err != nil {
			errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
			return
		}
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
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}	[get]
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
	path := filepath.Join(cmsConfig.ContentDir, collection, entry)

	blob, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errReposGetBlob().Status(resp.StatusCode).Log(r, err).Json(w)
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
	path := filepath.Join(cmsConfig.ContentDir, collection, entry)
	commitMessage := fmt.Sprintf("delete(%s): %s", collection, entry)

	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, nil, commitMessage)
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// components

// @Summary		Get components
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		sandpack		query	string	false	"response in sandpack format (true, false, 0 or 1)"
// @Success		200	{object}	map[string]string
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/components	[get]
// @Security	bearerToken
func getComponents(w http.ResponseWriter, r *http.Request) {
	cres, cached := cache.Get(r.URL.String())
	if cached {
		rawResponse(w, http.StatusOK, cres)
		return
	}

	ctx := r.Context()
	accessToken := accessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	sandpack, _ := strconv.ParseBool(r.FormValue("sandpack"))

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	rcs, resp, err := gh.GetContentsRecursive(ctx, accessToken, owner, repo, ref, cmsConfig.Components.EntryDir())
	if err != nil {
		errCmsGetComponents().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	componentsTree := make(map[string]string)
	for _, rc := range rcs {
		contentBytes, err := base64.StdEncoding.DecodeString(*rc.Content)
		if err != nil {
			errCmsGetComponents().Status(http.StatusInternalServerError).Log(r, err).Json(w)
			return
		}
		componentsTree[*rc.Path] = string(contentBytes)
	}

	if sandpack {
		files := make(map[string]interface{})
		for path, content := range componentsTree {
			files[path] = map[string]interface{}{
				"code": content,
			}
		}

		pkgJsonData, _, _ := gh.GetBlob(ctx, accessToken, owner, repo, ref, cms.PackageJSONFile)
		pkgJson := cms.ParsePackageJSON(pkgJsonData)

		jsonBytes := jsonResponse(w, http.StatusOK, map[string]interface{}{
			"files": files,
			"entry": cmsConfig.Components.Entry,
			"deps":  cms.SandpackResolveDeps(pkgJson, cmsConfig.Components.Dependencies),
		})

		cache.Set(r.URL.String(), jsonBytes)
		return
	}

	jsonResponse(w, http.StatusOK, componentsTree)
}
