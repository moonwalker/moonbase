package cli

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/moonbase/internal/runtime"
)

var RootCmd = &cobra.Command{
	Use:               runtime.Name,
	Short:             "Git-based headless CMS API",
	Version:           runtime.ShortRev(),
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Run() error {
	return RootCmd.Execute()
}
