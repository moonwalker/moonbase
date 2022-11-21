package cli

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/internal/version"
)

var RootCmd = &cobra.Command{
	Use:               version.Name,
	Short:             "Git-based headless CMS API",
	Version:           version.ShortRev(),
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Run() error {
	return RootCmd.Execute()
}
