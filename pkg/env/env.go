package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() {
	// .env (default)
	godotenv.Load()

	// .env.local # local user specific (git ignored)
	godotenv.Overload(".env.local")
}

func Getenv(key string, def string) string {
	s := os.Getenv(key)
	if len(s) == 0 {
		return def // return default
	}
	return s
}

func Int(key string, def int) int {
	i, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return def
	}
	return i
}
