package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/internal/gh"
)

type listItem struct {
	Name string `json:"name"`
}

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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	repos := make([]listItem, 0)
	for _, gr := range grs {
		repos = append(repos, listItem{
			Name: *gr.Name,
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := make([]listItem, 0)
	for _, br := range branches {
		bs = append(bs, listItem{
			Name: *br.Name,
		})
	}

	jsonResponse(w, http.StatusOK, bs)
}
