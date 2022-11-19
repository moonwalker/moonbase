package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
)

const respHtml = `
<html>
<head>
<head/>
<body>
<script>
window.onload = function() {
	var token, user;
	var match = document.cookie.match(new RegExp('(^| )gh_token=([^;]+)'));
  	if (match) token = match[2];
	match = document.cookie.match(new RegExp('(^| )artms_user=([^;]+)'));
  	if (match) user = match[2];
	window.opener.postMessage({ gh_token: token, artms_user: user }, '*');
	window.close();
};
</script>
</body
</html>`

const (
	oauthStateSeparator = "|"
)

var (
	ghScopes         = []string{"user:email", "read:org"}
	oauthStateSecret = xid.New().String()
)

type TokenData struct {
	Login       string `json:"login"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Login *string `json:"login,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Image *string `json:"image,omitempty"`
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
	githubConfig := githubConfig()

	secret, returnURL := decodeState(r)
	if secret != oauthStateSecret {
		httpError(w, -1, "invalid oauth state secret", fmt.Errorf("expected: %s, actual: %s", oauthStateSecret, secret))
		return
	}

	code := r.FormValue("code")

	// IDEA: here we can return with the 'code' (encrypted) to the client
	// and the client can initiate the oauth exchange step
	// so we won't need popup or iframe solutions
	println(returnURLWithCode(returnURL, code))

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
	// json.NewEncoder(w).Encode(User{
	// 	Login: ghUser.Login,
	// 	Name:  ghUser.Name,
	// 	Email: ghUser.Email,
	// 	Image: ghUser.AvatarURL,
	// })

	u, _ := json.Marshal(User{
		Login: ghUser.Login,
		Name:  ghUser.Name,
		Email: ghUser.Email,
		Image: ghUser.AvatarURL,
	})
	http.SetCookie(w, &http.Cookie{Name: "artms_user", Value: base64.StdEncoding.EncodeToString(u), Path: "/"})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, respHtml)
}

func encryptAccessToken(user *github.User, accessToken string) (string, error) {
	payload := &TokenData{Login: *user.Login, AccessToken: accessToken}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, data, 30)
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

	u := fmt.Sprintf("%s?code=%s", returnURL, codeJWT)
	return u, nil
}
