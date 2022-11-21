package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func httpError(w http.ResponseWriter, code int, msg string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	var e string
	if err != nil {
		e = err.Error()
	}
	json.NewEncoder(w).Encode(Error{code, msg, e})
}
