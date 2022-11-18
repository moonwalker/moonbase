package cli

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/app"
)

var (
	defaultPort = 8080
	port        int
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Moonbase server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		defaultPort = p
	}
	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", defaultPort, "HTTP port")
	RootCmd.AddCommand(serverCmd)
	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {
	srv := app.NewServer(&app.Options{Port: port})
	return srv.Listen()
}
