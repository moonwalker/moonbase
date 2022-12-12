package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/xid"

	"github.com/moonwalker/moonbase/internal/log"
)

// API uses conventional HTTP response codes to indicate the success or failure of an API request.
// Codes in the `2xx` range indicate success.
// Codes in the `4xx` range indicate an error that failed given the information provided (e.g., a required parameter was omitted, etc.).
// Codes in the `5xx` range indicate an error with servers (these are rare).

var (
	// common
	errNotFound   = errf(404, "err_http_001", "requested resource not found")
	errJsonDecode = errf(400, "err_json_001", "failed to decode request body")
	// auth
	errAuthNoToken     = errf(401, "err_auth_001", "auth token not provided")
	errAuthBadToken    = errf(401, "err_auth_002", "invalid auth token")
	errAuthBadClaims   = errf(400, "err_auth_003", "invalid auth claims")
	errAuthEncState    = errf(400, "err_auth_004", "failed to encode oauth state")
	errAuthBadSecret   = errf(400, "err_auth_005", "invalid oauth state secret")
	errAuthEncRetURL   = errf(400, "err_auth_006", "failed to encrypt return url with oauth exchange code")
	errAuthCodeMissing = errf(400, "err_auth_007", "auth code missing")
	errAuthDecOAuth    = errf(400, "err_auth_008", "failed to decrypt oauth exchange code")
	errAuthExchange    = errf(400, "err_auth_009", "oauth exchange failed")
	errAuthGetUser     = errf(400, "err_auth_010", "failed to get user")
	errAuthEncToken    = errf(400, "err_auth_011", "failed to encrypt token")
	// repos
	errReposGet         = errf(404, "err_repos_001", "failed to get repositories")
	errReposGetBranches = errf(404, "err_repos_002", "failed to get branches")
	errReposGetTree     = errf(404, "err_repos_003", "failed to get tree")
	errReposGetBlob     = errf(404, "err_repos_004", "failed to get blob")
	errReposCommitBlob  = errf(400, "err_repos_005", "failed to commit blob")
	errReposDeleteBlob  = errf(400, "err_repos_006", "failed to delete blob")
	// cms
	errCmsGetCommits       = errf(404, "err_cms_001", "failed to get commits")
	errCmsDeleteFolder     = errf(400, "err_cms_002", "failed to delete folder")
	errCmsSchemaValidation = errf(400, "err_cms_003", "schema validation failed")
	errCmsSchemaGeneration = errf(400, "err_cms_004", "schema generation failed")
	errCmsGetComponents    = errf(404, "err_cms_005", "failed to get components")
	errCmsParseBlob        = errf(400, "err_cms_006", "failed to parse blob")
)

type errorData struct {
	ID         string   `json:"id"`
	StatusCode int      `json:"statusCode"`
	StatusText string   `json:"statusText"`
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	Detailed   []string `json:"details,omitempty"`
}

func errf(statusCode int, code, message string) func() *errorData {
	return func() *errorData {
		id := xid.New().String()
		statusText := http.StatusText(statusCode)
		return &errorData{id, statusCode, statusText, code, message, nil}
	}
}

func (e *errorData) Details(a ...string) *errorData {
	for _, arg := range a {
		e.Detailed = append(e.Detailed, arg)
	}
	return e
}

func (e *errorData) Status(statusCode int) *errorData {
	e.StatusCode = statusCode
	return e
}

func (e *errorData) Json(w http.ResponseWriter) *errorData {
	jsonResponse(w, e.StatusCode, e)
	return e
}

func (e *errorData) Log(r *http.Request, err error) *errorData {
	reqid := middleware.GetReqID(r.Context())
	log.Error(err).
		Str("reqid", reqid).
		Str("id", e.ID).
		Int("status", e.StatusCode).
		Str("code", e.Code).
		Msg(e.Message)
	return e
}
