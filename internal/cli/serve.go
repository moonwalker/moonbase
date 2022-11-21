package cli

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/internal/env"
	"github.com/moonwalker/moonbase/internal/log"
	"github.com/moonwalker/moonbase/internal/server"
)

var (
	httpPort int
	serveCmd = &cobra.Command{
		Use:          "serve",
		Short:        "Run the Moonbase server",
		RunE:         serveCmdRun,
		SilenceUsage: true,
	}
)

func init() {
	serveCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", env.Port(8080), "HTTP port")
	RootCmd.AddCommand(serveCmd)
	RootCmd.RunE = serveCmdRun
}

func serveCmdRun(command *cobra.Command, args []string) error {
	log.Info().Int("port", httpPort).Msg("starting")
	srv := server.NewServer(&server.Options{Port: httpPort})
	return srv.Listen()
}
