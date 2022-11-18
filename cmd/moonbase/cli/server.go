package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/app"
	"github.com/moonwalker/moonbase/pkg/env"
)

var (
	port    int
	JWT_KEY = []byte(os.Getenv("JWT_KEY"))
	JWE_KEY = []byte(os.Getenv("JWE_KEY"))
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Moonbase server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", env.Int("PORT", 8080), "HTTP port")
	RootCmd.AddCommand(serverCmd)
	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {
	srv := app.NewServer(&app.Options{Port: port})
	return srv.Listen()
}
