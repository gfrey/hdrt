package server

import (
	"log"
	"net/http"
	"os"
	"strings"
)

var logger = log.New(os.Stderr, "", 0)

type Server struct {
}

func (r *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if !strings.Contains(port, ":") {
		port = "127.0.0.1:" + port
	}

	logger.Printf("running port %q", port)

	return http.ListenAndServe(port, Router())
}
