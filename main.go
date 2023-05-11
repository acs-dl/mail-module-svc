package main

import (
	"os"

	"github.com/acs-dl/mail-module-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
