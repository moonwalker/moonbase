package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v48/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
)

type TokenData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

var (
	githubConfig = &oauth2.Config{
		Scopes:   []string{"user:email", "read:org"},
		Endpoint: githuboauth.Endpoint,
	}
	oauthStateString = xid.New().String()
)

func githubAuth(w http.ResponseWriter, r *http.Request) {
	url := githubConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		httpError(w, -1, "invalid oauth state", fmt.Errorf("expected: %s, actual: %s", oauthStateString, state))
		return
	}

	code := r.FormValue("code")
	token, err := githubConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		httpError(w, -1, "oauth exchange failed", err)
		return
	}

	oauthClient := githubConfig.Client(oauth2.NoContext, token)
	ghClient := github.NewClient(oauthClient)
	ghUser, _, err := ghClient.Users.Get(context.Background(), "")
	if err != nil {
		httpError(w, -1, "github client failed to get user", err)
		return
	}

	et, err := encryptAccessToken(ghUser)
	if err != nil {
		httpError(w, -1, "failed to encrypt token", err)
		return
	}

	json.NewEncoder(w).Encode(User{Name: *ghUser.Name, Email: *ghUser.Email, Token: et})
}

func encryptAccessToken(user *github.User) (string, error) {
	payload := &TokenData{Name: *user.Name, Email: *user.Email}
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, payload, 1)
	if err != nil {
		return "", err
	}

	return te, nil
}
