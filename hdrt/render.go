package main

import (
	"os"

	"github.com/gfrey/hdrt"
	"github.com/gfrey/hdrt/ui"
)

type renderCmd struct {
	SzenePath string `cli:"arg desc='scene description to be rendered; stdin if none given'"`
}

func (cmd *renderCmd) Run() error {
	var err error
	var wrld *hdrt.World
	switch cmd.SzenePath {
	case "":
		wrld, err = hdrt.LoadWorldFromReader(os.Stdin)
		if err != nil {
			return err
		}
	default:
		wrld, err = hdrt.LoadWorldFromFile(cmd.SzenePath)
		if err != nil {
			return err
		}
	}

	return ui.Run(wrld)
}