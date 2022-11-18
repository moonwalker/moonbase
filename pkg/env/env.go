package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// env keys
const (
	PORT    = "PORT"
	JWT_KEY = "JWT_KEY"
	JWE_KEY = "JWE_KEY"
)

func Load() {
	// .env (default)
	godotenv.Load()

	// .env.local # local user specific (git ignored)
	godotenv.Overload(".env.local")
}

func Port(def int) int {
	return getint(PORT, def)
}

func JwtKey() string {
	return get(JWT_KEY, "")
}

func JweKey() string {
	return get(JWE_KEY, "")
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
