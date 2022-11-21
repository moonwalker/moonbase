package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v48/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/internal/env"
	"github.com/moonwalker/moonbase/internal/jwt"
)

// testing the flow:
// http://localhost:8080/login/github?returnURL=/login/github/authenticate

const (
	userCtxKey          = "user-token"
	oauthStateSep       = "|"
	retUrlCodeQuery int = 0
	retUrlCodePath      = 1
)

var (
	ghScopes         = []string{"user:email", "read:org", "repo"}
	oauthStateSecret = xid.New().String()
)

type User struct {
	Login *string `json:"login,omitempty"`
	Image *string `json:"image,omitempty"`
	Token string  `json:"token"`
}

func githubConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:       ghScopes,
		Endpoint:     githuboauth.Endpoint,
		ClientID:     env.GithubClientID,
		ClientSecret: env.GithubClientSecret,
	}
}

func githubAuth(w http.ResponseWriter, r *http.Request) {
	state, err := encodeState(r)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to encode oauth state", err)
		return
	}

	url := githubConfig().AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	secret, returnURL := decodeState(r)
	if secret != oauthStateSecret {
		httpError(w, http.StatusInternalServerError, "invalid oauth state secret", fmt.Errorf("expected: %s, actual: %s", oauthStateSecret, secret))
		return
	}

	code := r.FormValue("code")
	url, err := returnURLWithCode(returnURL, code, retUrlCodeQuery)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to encrypt return url with oauth exchange code", err)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	if code == "" {
		code = chi.URLParam(r, "code")
		if code == "" {
			httpError(w, http.StatusInternalServerError, "auth code missing", nil)
			return
		}
	}

	decoded, err := decryptExchangeCode(code)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to decrypt oauth exchange code", err)
		return
	}

	githubConfig := githubConfig()
	token, err := githubConfig.Exchange(oauth2.NoContext, decoded)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "oauth exchange failed", err)
		return
	}

	oauthClient := githubConfig.Client(oauth2.NoContext, token)
	ghClient := github.NewClient(oauthClient)
	ghUser, _, err := ghClient.Users.Get(context.Background(), "")
	if err != nil {
		httpError(w, http.StatusInternalServerError, "github client failed to get user", err)
		return
	}

	et, err := encryptAccessToken(token.AccessToken)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to encrypt token", err)
		return
	}

	usr := &User{
		Login: ghUser.Login,
		Image: ghUser.AvatarURL,
		Token: et,
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	json.NewEncoder(w).Encode(usr)
}

func encryptAccessToken(accessToken string) (string, error) {
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(accessToken) /*data*/, 30)
	if err != nil {
		return "", err
	}

	return te, nil
}

func encodeState(r *http.Request) (string, error) {
	returnURL := r.FormValue("returnURL")

	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}

	if !u.IsAbs() {
		u, err = url.Parse(r.Referer())
		u.Path = returnURL
	}

	state := fmt.Sprintf("%s%s%s", oauthStateSecret, oauthStateSep, u)
	return base64.URLEncoding.EncodeToString([]byte(state)), nil
}

func decodeState(r *http.Request) (string, string) {
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))
	parts := strings.Split(string(state), oauthStateSep)
	return parts[0], parts[1]
}

func returnURLWithCode(returnURL, code string, m int) (string, error) {
	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}

	codeJWT, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(code), 30)
	if err != nil {
		return "", err
	}

	switch {
	case m == retUrlCodeQuery:
		u.RawQuery = url.Values{
			"code": {codeJWT},
		}.Encode()
	case m == retUrlCodePath:
		u.Path = path.Join(u.Path, codeJWT)
	}

	return u.String(), nil
}

func decryptExchangeCode(code string) (string, error) {
	token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, code)
	if err != nil {
		return "", err
	}

	return string(token.Claims.(*jwt.AuthClaims).Data), nil
}

func withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// get token from authorization header
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}
		if len(tokenString) == 0 {
			httpError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), fmt.Errorf("no auth token"))
			return
		}

		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString)
		if err != nil {
			httpError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		authClaims, ok := token.Claims.(*jwt.AuthClaims)
		if !ok {
			httpError(w, http.StatusInternalServerError, "invalid auth claims type", nil)
			return
		}

		// add auth claims to context
		ctx := context.WithValue(r.Context(), userCtxKey, string(authClaims.Data))

		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}