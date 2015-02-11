package main

import "github.com/dynport/dgtk/cli"

import "github.com/gfrey/hdrt/server"

func CLIRouter() *cli.Router {
	r := cli.NewRouter()

	r.Register("server", &server.Server{}, "start the server")

	return r
}
