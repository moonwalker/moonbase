package cli

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/app"
	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/version"
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
	log.Printf("moonbase version %s", version.ShortRev())
	srv := app.NewServer(&app.Options{Port: port})
	return srv.Listen()
}
