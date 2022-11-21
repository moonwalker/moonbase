package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v48/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
)

// testing the flow:
// http://localhost:8080/login/github?returnURL=/login/github/authenticate

const (
	oauthStateSeparator = "|"
	USER_CTX_KEY        = "user-token"
)

var (
	ghScopes         = []string{"user:email", "read:org", "repo"}
	oauthStateSecret = xid.New().String()
)

type User struct {
	Login *string `json:"login,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Image *string `json:"image,omitempty"`
	Token string  `json:"token,omitempty"`
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
	state := encodeState(r)
	url := githubConfig().AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	secret, returnURL := decodeState(r)
	if secret != oauthStateSecret {
		httpError(w, -1, "invalid oauth state secret", fmt.Errorf("expected: %s, actual: %s", oauthStateSecret, secret))
		return
	}

	code := r.FormValue("code")
	url, err := returnURLWithCode(returnURL, code)
	if err != nil {
		httpError(w, -1, "failed to encrypt return url with oauth exchange code", err)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authenticateCallback(w http.ResponseWriter, r *http.Request) {
	githubConfig := githubConfig()
	code := chi.URLParam(r, "code")
	if code == "" {
		httpError(w, -1, "auth code missing", nil)
		return
	}

	decoded, err := decryptExchangeCode(code)
	if err != nil {
		httpError(w, -1, "failed to decrypt oauth exchange code", err)
		return
	}

	token, err := githubConfig.Exchange(oauth2.NoContext, decoded)
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

	et, err := encryptAccessToken(token.AccessToken)
	if err != nil {
		httpError(w, -1, "failed to encrypt token", err)
		return
	}

	usr := &User{
		Login: ghUser.Login,
		Name:  ghUser.Name,
		Email: ghUser.Email,
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

func encodeState(r *http.Request) string {
	returnURL := r.FormValue("returnURL")
	if returnURL == "" {
		returnURL = r.Referer()
	}
	state := fmt.Sprintf("%s%s%s", oauthStateSecret, oauthStateSeparator, returnURL)
	return base64.URLEncoding.EncodeToString([]byte(state))
}

func decodeState(r *http.Request) (string, string) {
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))
	parts := strings.Split(string(state), oauthStateSeparator)
	return parts[0], parts[1]
}

func returnURLWithCode(returnURL, code string) (string, error) {
	codeJWT, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(code), 30)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("%s/%s", returnURL, codeJWT)
	return u, nil
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
		ctx := context.WithValue(r.Context(), USER_CTX_KEY, string(authClaims.Data))

		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
