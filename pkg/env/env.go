package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	JwtKey string
	JweKey string
)

func Load() {
	// .env (default)
	godotenv.Load()

	// .env.local # local user specific (git ignored)
	godotenv.Overload(".env.local")

	// set vars
	JwtKey = os.Getenv("JWT_KEY")
	JweKey = os.Getenv("JWE_KEY")
}

func Port(def int) int {
	return getint("PORT", def)
}

// private functions

func get(key string, def string) string {
	s := os.Getenv(key)
	if len(s) == 0 {
		return def // return default
	}
	return s
}

func getint(key string, def int) int {
	i, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return def
	}
	return i
}
