package gh

import (
	"context"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/internal/env"
)

var (
	ghScopes = []string{"user:email", "read:org", "repo"}
)

func ghConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:       ghScopes,
		Endpoint:     githuboauth.Endpoint,
		ClientID:     env.GithubClientID,
		ClientSecret: env.GithubClientSecret,
	}
}

func ghClient(ctx context.Context, accessToken string) *github.Client {
	oauthClient := ghConfig().Client(ctx, &oauth2.Token{AccessToken: accessToken})
	return github.NewClient(oauthClient)
}

func AuthCodeURL(state string) string {
	return ghConfig().AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func Exchange(code string) (string, error) {
	t, err := ghConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}

func GetUser(ctx context.Context, accessToken string) (*github.User, error) {
	ghClient := ghClient(ctx, accessToken)
	user, _, err := ghClient.Users.Get(ctx, "")
	return user, err
}

func ListRepositories(ctx context.Context, accessToken string, page, perPage int, sort, direction string) ([]*github.Repository, int, error) {
	repos, resp, err := ghClient(ctx, accessToken).Repositories.List(ctx, "", &github.RepositoryListOptions{
		Sort:      sort,
		Direction: direction,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	})
	return repos, resp.LastPage, err
}

func ListBranches(ctx context.Context, accessToken string, owner, repo string) ([]*github.Branch, error) {
	branches, _, err := ghClient(ctx, accessToken).Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{})
	return branches, err
}
