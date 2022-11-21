package main

import (
	"os"

	"github.com/moonwalker/moonbase/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		os.Exit(1)
	}
}
