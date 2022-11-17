package cli

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type Command = cobra.Command

var RootCmd = &cobra.Command{
	Use:   "moonbase",
	Short: "Git-based headless CMS API",
}

func init() {
	// .env (default)
	godotenv.Load()

	// .env.local # local user specific (git ignored)
	godotenv.Overload(".env.local")
}

func Run() error {
	return RootCmd.Execute()
}
