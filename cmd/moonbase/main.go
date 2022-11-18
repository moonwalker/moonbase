package main

import (
	"os"

	"github.com/moonwalker/moonbase/cmd/moonbase/cli"
	"github.com/moonwalker/moonbase/pkg/env"
)

func init() {
	env.Load()
}

func main() {
	if err := cli.Run(); err != nil {
		os.Exit(1)
	}
}
