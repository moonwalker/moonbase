package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/internal/env"
	"github.com/moonwalker/moonbase/internal/log"
	"github.com/moonwalker/moonbase/internal/runtime"
	"github.com/moonwalker/moonbase/internal/server"
)

var (
	httpPort int
	serveCmd = &cobra.Command{
		Use:          "serve",
		Short:        "Run moonbase server",
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
	k := "port"
	v := fmt.Sprintf("%d", httpPort)

	if runtime.IsDev() {
		k = "docs"
		v = fmt.Sprintf("http://localhost:%d/docs", httpPort)
	}

	log.Info().Str(k, v).Msg("running")
	return server.Listen(httpPort)
}
