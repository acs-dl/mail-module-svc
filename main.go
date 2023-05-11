package main

import (
	"os"

	"gitlab.com/distributed_lab/acs/mail-module/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
