package api

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	JWT_KEY = []byte(os.Getenv("JWT_KEY"))
	JWE_KEY = []byte(os.Getenv("JWE_KEY"))
)

func Debug(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Time: %s", time.Now().UTC())
}
