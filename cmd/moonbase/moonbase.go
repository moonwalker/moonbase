package main

import (
	"os"

	"github.com/moonwalker/moonbase/cmd/moonbase/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		os.Exit(1)
	}
}
