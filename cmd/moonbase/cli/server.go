package cli

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/app"
	"github.com/moonwalker/moonbase/pkg/env"
)

var (
	port int
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Moonbase server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", env.Port(8080), "HTTP port")
	RootCmd.AddCommand(serverCmd)
	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {
	srv := app.NewServer(&app.Options{Port: port})
	return srv.Listen()
}
