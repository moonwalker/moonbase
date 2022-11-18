package cli

import (
	"github.com/moonwalker/moonbase/pkg/env"

	"github.com/spf13/cobra"
)

type Command = cobra.Command

var RootCmd = &cobra.Command{
	Use:   "moonbase",
	Short: "Git-based headless CMS API",
}

func init() {
	env.Load()
}

func Run() error {
	return RootCmd.Execute()
}
