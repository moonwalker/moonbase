package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v48/github"
	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
	"golang.org/x/oauth2"
)

func githubAuth(w http.ResponseWriter, r *http.Request) {
	url := githubConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		redirectWithError(w, r, "/", "invalid oauth state", fmt.Errorf("expected: %s, actual: %s", oauthStateString, state))
		return
	}

	code := r.FormValue("code")
	token, err := githubConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		redirectWithError(w, r, "/", "oauth exchange failed", err)
		return
	}

	oauthClient := githubConfig.Client(oauth2.NoContext, token)
	ghClient := github.NewClient(oauthClient)
	ghUser, _, err := ghClient.Users.Get(context.Background(), "")
	if err != nil {
		redirectWithError(w, r, "/", "github client failed to get user", err)
		return
	}

	et, err := encryptAccessToken(ghUser)
	if err != nil {
		redirectWithError(w, r, "/", "failed to encrypt token", err)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "gh_token", Value: et, Path: "/"})
	json.NewEncoder(w).Encode(User{Name: *ghUser.Name, Email: *ghUser.Email})
}

func redirectWithError(w http.ResponseWriter, r *http.Request, url string, msg string, err error) {
	http.SetCookie(w, &http.Cookie{Name: "FLASH_ERROR", Value: base64.URLEncoding.EncodeToString([]byte(msg)), Path: url})
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func encryptAccessToken(user *github.User) (string, error) {
	payload := &TokenData{Name: *user.Name, Email: *user.Email}
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, payload, 1)
	if err != nil {
		return "", err
	}

	return te, nil
}
