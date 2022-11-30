package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func rawResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func jsonResponse(w http.ResponseWriter, statusCode int, v any) []byte {
	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	res := buf.Bytes()
	rawResponse(w, statusCode, res)
	return res
}
