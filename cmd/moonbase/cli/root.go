package cli

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/pkg/version"
)

type Command = cobra.Command

var RootCmd = &cobra.Command{
	Use:     "moonbase",
	Short:   "Git-based headless CMS API",
	Version: version.ShortRev(),
}

func Run() error {
	return RootCmd.Execute()
}
