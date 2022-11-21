package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type listItem struct {
	Name string `json:"name"`
}

func createClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func getRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)

	githubClient := createClient(accessToken)
	grs, _, err := githubClient.Repositories.List(ctx, "", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	})
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

	json.NewEncoder(w).Encode(repos)
}

func getBranches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessToken := ctx.Value(userCtxKey).(string)
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	githubClient := createClient(accessToken)
	branches, _, err := githubClient.Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{})
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

	json.NewEncoder(w).Encode(bs)
}
