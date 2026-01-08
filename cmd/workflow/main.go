package main

import (
	"os"

	"github.com/zevwings/workflow/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
