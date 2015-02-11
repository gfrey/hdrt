package server

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

type Server struct {
}

func (r *Server) Run() error {
	logger.Printf("running")
	return nil
}
