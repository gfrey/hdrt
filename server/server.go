package server

import (
	"fmt"
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

	renderDir := os.TempDir() + "/hdrt"
	err := prepareRenderDir(renderDir)
	if err != nil {
		return err
	}

	return http.ListenAndServe(port, Router(renderDir))
}

func prepareRenderDir(renderDir string) error {
	_, err := os.Stat(renderDir)

	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Unknwo error when getting stats: %s", err.Error())
		}
		logger.Printf("Will create renderDir %s", renderDir)
	} else {
		// does exist, delete
		err := os.RemoveAll(renderDir)
		if err != nil {
			return fmt.Errorf("error removing file %s (%s)", renderDir, err.Error())
		}
		logger.Printf("deleted and will recreate renderDir %s", renderDir)
	}

	err = os.Mkdir(renderDir, 0777)
	if err != nil {
		return fmt.Errorf("error when creating dir: %s", err.Error())
	}

	return nil
}
