package api

import (
	"fmt"
	"net/http"
	"time"
)

func Debug(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Time: %s", time.Now().UTC())
}
