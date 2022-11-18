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
	Login       string `json:"login"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
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

	et, err := encryptAccessToken(ghUser, token.AccessToken)
	if err != nil {
		httpError(w, -1, "failed to encrypt token", err)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "gh_token", Value: et, Path: "/"})
	json.NewEncoder(w).Encode(User{
		Name:  *ghUser.Name,
		Email: *ghUser.Email,
		Image: *ghUser.AvatarURL,
	})
}

func encryptAccessToken(user *github.User, accessToken string) (string, error) {
	payload := &TokenData{Login: *user.Login, AccessToken: accessToken}
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, payload, 30)
	if err != nil {
		return "", err
	}

	return te, nil
}

func returnWithError(w http.ResponseWriter, code int, msg string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(Error{code, msg, err.Error()})
}
