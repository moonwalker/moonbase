package gh

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"gopkg.in/yaml.v2"

	"github.com/moonwalker/moonbase/internal/env"
)

type cmsConfig struct {
	Content struct {
		Dir string `yaml:"dir"`
	} `yaml:"content"`
}

var (
	ghScopes = []string{"user:email", "read:org", "repo"}
)

const (
	configYamlPath = "moonbase.yaml"
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
	user, _, err := ghClient(ctx, accessToken).Users.Get(ctx, "")
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

func GetTree(ctx context.Context, accessToken string, owner string, repo string, branch string, sha string) (*github.Tree, error) {
	githubClient := ghClient(ctx, accessToken)

	// get branch's base ref
	if len(sha) == 0 {
		baseRef, _, err := githubClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
		if err != nil {
			return nil, err
		}
		sha = *baseRef.Object.SHA

		contentDir, err := getContentDir(ctx, githubClient, owner, repo, branch)
		if contentDir != "" {
			contentTree, _, err := githubClient.Git.GetTree(ctx, owner, repo, sha, false)
			if err != nil {
				return nil, err
			}

			for _, cti := range contentTree.Entries {
				if *cti.Type == "tree" && *cti.Path == contentDir {
					sha = *cti.SHA
					break
				}
			}
		}
	}

	// get tree
	tree, _, err := githubClient.Git.GetTree(ctx, owner, repo, sha, false)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func GetBlobByPath(ctx context.Context, accessToken string, owner string, repo string, branch, path string) ([]byte, error) {
	githubClient := ghClient(ctx, accessToken)
	blob, err := getBlobByPath(ctx, githubClient, owner, repo, branch, path)
	if err != nil {
		return nil, err
	}
	return blob, nil
}

// helpers
func getContentDir(ctx context.Context, githubClient *github.Client, owner string, repo string, branch string) (string, error) {
	cfg, err := getCmsConfig(ctx, githubClient, owner, repo, branch)
	if err != nil {
		return "", err
	}
	return cfg.Content.Dir, nil
}

func getCmsConfig(ctx context.Context, githubClient *github.Client, owner string, repo string, branch string) (*cmsConfig, error) {
	config, err := getBlobByPath(ctx, githubClient, owner, repo, branch, configYamlPath)
	if err != nil {
		return nil, err
	}

	cfg := &cmsConfig{}
	err = yaml.Unmarshal(config, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getBlobByPath(ctx context.Context, githubClient *github.Client, owner string, repo string, branch string, path string) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("No path provided")
	}

	fc, _, _, err := githubClient.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}

	decodedBlob, err := base64.StdEncoding.DecodeString(*fc.Content)
	if err != nil {
		return nil, err
	}

	return decodedBlob, nil
}
