package main

import (
	"log"
	"os"

	"github.com/dynport/dgtk/cli"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	router := CLIRouter()

	switch err := router.RunWithArgs(); err {
	case nil, cli.ErrorHelpRequested, cli.ErrorNoRoute:
		// ignore
		return
	default:
		logger.Fatal(err)
	}
}
