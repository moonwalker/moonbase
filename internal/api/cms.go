package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"

	"github.com/moonwalker/moonbase/pkg/content"

	"github.com/moonwalker/moonbase/internal/cache"
	"github.com/moonwalker/moonbase/internal/cms"
	gh "github.com/moonwalker/moonbase/pkg/github"
)

type collectionPayload struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type entryPayload struct {
	Login      string `json:"login"`
	Name       string `json:"name"`
	Contents   string `json:"contents"`
	SaveSchema bool   `json:"save_schema"`
}

type commitEntry struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type ComponentsTree map[string]string
type ComponentsTreeSha string

type localizedEntry struct {
	Name    string                     `json:"name"`
	Type    string                     `json:"type"`
	Content *content.MergedContentData `json:"content"`
	Schema  content.Schema             `json:"schema,omitempty"`
}

type entryItem struct {
	Name    string               `json:"name"`
	Content *content.ContentData `json:"content"`
	Schema  content.Schema       `json:"schema,omitempty"`
}

type collectionGroup struct {
	Label  string   `json:"label"`
	Models []string `json:"models"`
}

type entryResponse struct {
}

type queryParams struct {
	Page    int64  `url:"page"`
	Limit   int64  `url:"limit"`
	Include int64  `url:"include"`
	Filter  string `url:"filter"`
	Order   string `url:"order"`
}

type pagination struct {
	CurrentPage  int64  `json:"currentPage"`
	TotalCount   int64  `json:"totalCount"`
	NextPage     *int64 `json:"nextPage,omitempty"`
	PreviousPage *int64 `json:"prevPage,omitempty"`
	PageCount    *int64 `json:"pageCount,omitempty"`
}

type filters map[string]string
type orderBy map[string]int

type listResponse struct {
	Data       []*content.ContentData `json:"data"`
	Schema     *content.Schema        `json:"schema"`
	Pagination *pagination            `json:"pagination,omitempty"`
	Filter     *filters               `json:"filter,omitempty"`
	Order      *orderBy               `json:"order,omitempty"`
}

var (
	shaCache = cache.NewGeneric[ComponentsTreeSha](30 * time.Minute)
)

func commitMessage(collection, method, name string) string {
	return fmt.Sprintf("feat(%s): %s %s", collection, method, name)
}

// config

func getConfig(ctx context.Context, accessToken string, owner string, repo string, ref string) *cms.Config {
	data, _, _ := gh.GetBlob(ctx, accessToken, owner, repo, ref, cms.ConfigPath)
	return cms.ParseConfig(data)
}

// schema

func getSchema(ctx context.Context, accessToken string, owner string, repo string, ref string, collection string, workdir string) *cms.Schema {
	p := filepath.Join(workdir, collection, content.JsonSchemaName)
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
	accessToken := gh.AccessTokenFromContext(ctx)

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
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, cmsConfig.WorkDir)
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
	accessToken := gh.AccessTokenFromContext(ctx)

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
	path := filepath.Join(cmsConfig.WorkDir, collectionName, content.JsonSchemaName)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	emptyContent := fmt.Sprintf(`{"id":"%s","name":"%s","displayField":"","fields":[],"createdAt":"%s","createdBy":"%s","updatedAt":"%s","updatedBy":"%s","version":0}`, collection.Name, cases.Title(language.Und, cases.NoLower).String(collection.Name), now, collection.Login, now, collection.Login)

	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &emptyContent, commitMessage("content", "create", collectionName))
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
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

	collectionName := slug.Make(collection)
	path := filepath.Join(cmsConfig.WorkDir, collectionName)
	if path == cmsConfig.WorkDir {
		m := "missing collection name"
		errCmsDeleteFolder().Details(m).Log(r, errors.New(m)).Json(w)
		return
	}

	resp, err := gh.DeleteFolder(ctx, accessToken, owner, repo, ref, path, commitMessage("content", "delete", collectionName))
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
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.WorkDir, collection)

	repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errReposGetTree().Status(resp.StatusCode).Log(r, err).Json(w)
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
	accessToken := gh.AccessTokenFromContext(ctx)

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
		return
	}

	entryData.Name = strings.ToLower(entryData.Name)

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	// if !entryData.SaveSchema {
	// 	schema := getSchema(ctx, accessToken, owner, repo, ref, collection, cmsConfig.WorkDir)
	// 	if schema.Available() {
	// 		err = schema.ValidateString(entryData.Contents)
	// 		if err != nil {
	// 			errCmsSchemaValidation().Log(r, err).Json(w)
	// 			return
	// 		}
	// 	}
	// }

	contentData := content.MergedContentData{}
	err = json.Unmarshal([]byte(entryData.Contents), &contentData)
	if err != nil {
		errCmsReadContent().Log(r, err).Json(w)
		return
	}
	now := time.Now().UTC()
	if len(contentData.Name) == 0 {
		contentData.Name = entryData.Name
		contentData.ID = entryData.Name
		contentData.CreatedAt = &now
		contentData.CreatedBy = entryData.Login
		contentData.Version = 1
		contentData.Status = "draft"
	} else {
		contentData.UpdatedAt = &now
		contentData.UpdatedBy = entryData.Login
		contentData.Version = contentData.Version + 1
		contentData.Status = "changed"
	}

	locales, statusCode, err := getLocales(ctx, accessToken, owner, repo, ref)
	if err != nil {
		errReposGetTree().Status(statusCode).Log(r, err).Json(w)
		return
	}

	items, err := cms.SeparateLocalisedContent(contentData, locales, cmsConfig.WorkDir, collection)
	if err != nil {
		errCmsSeparateLocalizedContent().Log(r, err).Json(w)
		return
	}

	resp, err := gh.CommitBlobs(ctx, accessToken, owner, repo, ref, items, commitMessage(collection, "create/update", entryData.Name))
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}
	// if ext == ".md" || ext == ".mdx" {
	// 	contentData, err = cms.JsonToMarkdown(contentData)
	// 	if err != nil {
	// 		errCmsParseMarkdown().Log(r, err).Json(w)
	// 		return
	// 	}
	// }

	// TODO: Uncomment to save to GH
	// resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &contentData, commitMessage)
	// if err != nil {
	// 	errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
	// 	return
	// }

	// if entryData.SaveSchema {
	// 	schema, err := cms.GenerateSchema(entryData.Name, entryData.Contents)
	// 	if err != nil {
	// 		errCmsSchemaGeneration().Log(r, err).Json(w)
	// 		return
	// 	}
	// 	schemaPath := filepath.Join(cmsConfig.WorkDir, collection, content.JsonSchemaName)
	// 	schemaCommitMessage := fmt.Sprintf("feat(%s): create/update %s", collection, content.JsonSchemaName)
	// 	resp, err = gh.CommitBlob(ctx, accessToken, owner, repo, ref, schemaPath, &schema, schemaCommitMessage)
	// 	if err != nil {
	// 		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
	// 		return
	// 	}
	// }

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
// @Success		200	{object}	map[string]interface{}
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}	[get]
// @Security	bearerToken
func getEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")
	entry := chi.URLParam(r, "entry")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	schemaPath := filepath.Join(cmsConfig.WorkDir, collection, content.JsonSchemaName)
	sc, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, schemaPath)
	if err != nil {
		errReposGetBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}
	cs := &content.Schema{}
	err = json.Unmarshal(sc, &cs)
	if err != nil {
		errCmsParseSchema().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	var mc *content.MergedContentData

	if entry != "_new" {
		// Get files in directory
		path := filepath.Join(cmsConfig.WorkDir, collection, entry)
		rc, resp, err := gh.GetAllLocaleContents(ctx, accessToken, owner, repo, ref, path)
		if err != nil {
			errReposGetBlob().Status(resp.StatusCode).Log(r, err).Json(w)
			return
		}

		mc, err = cms.MergeLocalisedContent(rc, *cs)
		if err != nil {
			errCmsMergeLocalizedContent().Log(r, err).Json(w)
			return
		}
	} else {
		locales, statusCode, err := getLocales(ctx, accessToken, owner, repo, ref)
		if err != nil {
			errReposGetTree().Status(statusCode).Log(r, err).Json(w)
			return
		}
		mc, err = cms.GetEmptyLocalisedContent(*cs, locales)
		if err != nil {
			errCmsMergeLocalizedContent().Log(r, err).Json(w)
			return
		}
	}

	data := &localizedEntry{Name: mc.Name, Type: mc.Type, Content: mc, Schema: *cs}
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
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")
	entry := chi.URLParam(r, "entry")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	path := filepath.Join(cmsConfig.WorkDir, collection, entry)
	resp, err := gh.DeleteFolder(ctx, accessToken, owner, repo, ref, path, commitMessage(collection, "delete", entry))
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// images

// @Summary		Upload image
// @Tags		cms
// @Accept		multipart/form-data
// @Param		owner			path		string		true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path		string		true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path		string		true	"git ref (branch, tag, sha)"
// @Param		payload			formData	file		true	"uploaded image"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/images	[post]
// @Security	bearerToken
func postImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	fileName := ""

	reader, err := r.MultipartReader()
	if err != nil {
		errCmsGetFormReader().Log(r, err).Json(w)
		return
	}

	var imgbytes []byte

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			errCmsReadForm().Log(r, err).Json(w)
			return
		}
		defer part.Close()
		fileBytes, err := io.ReadAll(part)
		if err != nil {
			errCmsReadContent().Log(r, err).Json(w)
			return
		}

		if fileName == "" {
			fileName = part.FileName()
		}

		imgbytes = append(imgbytes, fileBytes...)
	}
	imgbytes = bytes.Trim(imgbytes, "\xef\xbb\xbf")

	path := filepath.Join(cms.ImagesFolder, fileName)
	encoding := "base64"
	content := base64.StdEncoding.EncodeToString(imgbytes)

	blob, resp, err := gh.CreateBlob(ctx, accessToken, owner, repo, ref, &content, &encoding)
	if err != nil {
		errReposCreateBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}
	resp, err = gh.CommitBlobs(ctx, accessToken, owner, repo, ref, []gh.BlobEntry{
		{
			Path: path,
			SHA:  blob.SHA,
		}}, commitMessage("images", "upload", fileName))
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	data := &entryItem{Name: fileName}
	jsonResponse(w, http.StatusOK, data)
}

// settings

// @Summary		Get settings
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Success		200	{object}	[]entryPayload
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/settings	[get]
// @Security	bearerToken
func getSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, cms.SettingsFolder)
	if err != nil {
		errReposGetTree().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	items := make([]*entryPayload, 0)
	for _, rc := range repoContents {
		if *rc.Type == "file" {
			items = append(items, &entryPayload{
				Name: *rc.Name,
			})
		}
	}

	jsonResponse(w, http.StatusOK, items)
}

// @Summary		Get setting
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		setting			path	string	true	"setting"
// @Success		200	{object}	entryPayload
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/settings/{setting}	[get]
// @Security	bearerToken
func getSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	setting := chi.URLParam(r, "setting")

	name := setting + ".json"
	path := filepath.Join(cms.SettingsFolder, name)
	blob, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		e := errReposGetBlob()
		if resp != nil {
			e.Status(resp.StatusCode)
		}
		e.Log(r, err).Json(w)
		return
	}

	jsonResponse(w, http.StatusOK, &entryPayload{
		Name:     name,
		Contents: string(blob),
	})
}

// @Summary		Create/Update setting
// @Tags		cms
// @Accept		json
// @Param		owner			path	string			true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string			true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string			true	"git ref (branch, tag, sha)"
// @Param		collection		path	string			true	"collection"
// @Param		setting			path	string			true	"setting"
// @Param		payload			body	[]data			true	"setting payload"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/settings/{setting}	[post]
// @Security	bearerToken
func postSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	setting := chi.URLParam(r, "setting")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		errJsonDecode().Log(r, err).Json(w)
		return
	}

	path := filepath.Join(cms.SettingsFolder, strings.ToLower(setting)+".json")
	contents := string(b)
	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, &contents, commitMessage("settings", "create/update", setting))
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Delete setting
// @Tags		cms
// @Accept		json
// @Param		owner			path	string		true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string		true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string		true	"git ref (branch, tag, sha)"
// @Param		setting			path	string		true	"setting"
// @Success		200
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/settings/{setting}	[delete]
// @Security	bearerToken
func delSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	setting := chi.URLParam(r, "setting")

	path := filepath.Join(cms.SettingsFolder, setting+".json")
	resp, err := gh.CommitBlob(ctx, accessToken, owner, repo, ref, path, nil, commitMessage("settings", "delete", setting))
	if err != nil {
		errReposCommitBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Get reference
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		collection		path	string	true	"collection"
// @Param		id				path	string	true	"id"
// @Param		locale			path	string	true	"locale"
// @Success		200	{object}	map[string]interface{}
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/reference/{collection}/{id}/{locale}	[get]
// @Security	bearerToken
func getReference(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	collection := chi.URLParam(r, "collection")
	id := chi.URLParam(r, "id")
	locale := chi.URLParam(r, "locale")

	cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)
	schemaPath := filepath.Join(cmsConfig.WorkDir, collection, content.JsonSchemaName)
	sc, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, schemaPath)
	if err != nil {
		errReposGetBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}
	cs := &content.Schema{}
	err = json.Unmarshal(sc, &cs)
	if err != nil {
		errCmsParseSchema().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	// Get files in directory
	path := filepath.Join(cmsConfig.WorkDir, collection, id, locale+".json")
	rc, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		errReposGetBlob().Status(resp.StatusCode).Log(r, err).Json(w)
		return
	}

	contentData := content.ContentData{}
	err = json.Unmarshal(rc, &contentData)
	if err != nil {
		errCmsParseBlob().Status(http.StatusInternalServerError).Log(r, err).Json(w)
		return
	}

	jsonResponse(w, http.StatusOK, &entryItem{Schema: *cs, Content: &contentData})
}

// collection groups

// @Summary		Get collection groups
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Success		200	{object}	map[string]collectionGroup
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collectiongroups	[get]
// @Security	bearerToken
func getCollectionGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")

	path := filepath.Join(cms.SettingsFolder, "collectiongroups.json")
	blob, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		e := errReposGetBlob()
		if resp != nil {
			e.Status(resp.StatusCode)
		}
		e.Log(r, err).Json(w)
		return
	}

	groups := make(map[string]collectionGroup)
	err = json.Unmarshal(blob, &groups)
	if err != nil {
		errCmsParseBlob().Status(http.StatusInternalServerError).Log(r, err).Json(w)
		return
	}

	jsonResponse(w, http.StatusOK, groups)
}

// @Summary		Get collection group
// @Tags		cms
// @Accept		json
// @Produce		json
// @Param		owner			path	string	true	"the account owner of the repository (the name is not case sensitive)"
// @Param		repo			path	string	true	"the name of the repository (the name is not case sensitive)"
// @Param		ref				path	string	true	"git ref (branch, tag, sha)"
// @Param		group			path	string	true	"collection group"
// @Success		200	{object}	[]treeItem
// @Failure		500	{object}	errorData
// @Router		/cms/{owner}/{repo}/{ref}/collectiongroups/{group}	[get]
// @Security	bearerToken
func getCollectionGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := gh.AccessTokenFromContext(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	ref := chi.URLParam(r, "ref")
	group := chi.URLParam(r, "group")

	path := filepath.Join(cms.SettingsFolder, "collectiongroups.json")
	blob, resp, err := gh.GetBlob(ctx, accessToken, owner, repo, ref, path)
	if err != nil {
		e := errReposGetBlob()
		if resp != nil {
			e.Status(resp.StatusCode)
		}
		e.Log(r, err).Json(w)
		return
	}

	groups := make(map[string]collectionGroup)
	err = json.Unmarshal(blob, &groups)
	if err != nil {
		errCmsParseBlob().Status(http.StatusInternalServerError).Log(r, err).Json(w)
		return
	}

	for n, g := range groups {
		if n == group {
			collections := make(map[string]bool)
			for _, c := range g.Models {
				collections[c] = true
			}
			cmsConfig := getConfig(ctx, accessToken, owner, repo, ref)

			repoContents, resp, err := gh.GetTree(ctx, accessToken, owner, repo, ref, cmsConfig.WorkDir)
			if err != nil {
				errReposGetTree().Status(resp.StatusCode).Log(r, err).Json(w)
				return
			}

			treeItems := make([]*treeItem, 0)
			for _, rc := range repoContents {
				if *rc.Type == "dir" && collections[*rc.Name] {
					treeItems = append(treeItems, &treeItem{
						Name: rc.Name,
						Type: rc.Type,
						SHA:  rc.SHA,
					})
				}
			}

			jsonResponse(w, http.StatusOK, treeItems)
			return
		}
	}

	b := []byte{}
	jsonResponse(w, http.StatusOK, b)
}
