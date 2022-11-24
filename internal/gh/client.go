package gh

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

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

func GetTree(ctx context.Context, accessToken string, owner string, repo string, branch string, path string) ([]*github.RepositoryContent, error) {
	_, rc, _, err := ghClient(ctx, accessToken).Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: branch,
	})
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func GetBlob(ctx context.Context, accessToken string, owner string, repo string, ref, path string) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("path not provided")
	}

	fc, _, _, err := ghClient(ctx, accessToken).Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return nil, err
	}

	decodedBlob, err := base64.StdEncoding.DecodeString(*fc.Content)
	if err != nil {
		return nil, err
	}

	return decodedBlob, nil
}

func CommitBlob(ctx context.Context, accessToken string, owner string, repo string, ref string, path string, user string, email string, content string, commitMessage string) error {
	githubClient := ghClient(ctx, accessToken)

	reference, _, err := githubClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+ref)
	if err != nil {
		return err
	}

	tree, err := getCommitTree(ctx, githubClient, owner, repo, *reference.Object.SHA, path, content)
	if err != nil {
		return err
	}

	if err := pushCommit(ctx, githubClient, reference, tree, owner, repo, user, email, commitMessage); err != nil {
		return err
	}

	return nil
}

// getCommitTree generates the tree to commit based on the given files and the commit of the ref
func getCommitTree(ctx context.Context, githubClient *github.Client, owner string, repo string, sha string, path string, content string) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{
		{
			Path:    github.String(path),
			Type:    github.String("blob"),
			Content: github.String(content),
			Mode:    github.String("100644"),
		},
	}

	tree, _, err = githubClient.Git.CreateTree(ctx, owner, repo, sha, entries)
	return tree, err
}

// pushCommit creates the commit in the given reference using the given tree
func pushCommit(ctx context.Context, githubClient *github.Client, ref *github.Reference, tree *github.Tree, owner string, repo string, user string, email string, commitMessage string) (err error) {
	// Get the parent commit
	parent, _, err := githubClient.Repositories.GetCommit(ctx, owner, repo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &user, Email: &email}
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := githubClient.Git.CreateCommit(ctx, owner, repo, commit)

	if err != nil {
		return err
	}

	// Attach the commit to the branch
	ref.Object.SHA = newCommit.SHA
	_, _, err = githubClient.Git.UpdateRef(ctx, owner, repo, ref, false)
	return err
}
