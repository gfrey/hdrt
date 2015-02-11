package hdrt

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Renderer struct {
	wrld   *World
	evChan chan string
	abort  chan struct{}
}

func (r *Renderer) Render(renderDir string) (<-chan string, error) {
	if r.wrld == nil {
		return nil, fmt.Errorf("world not loaded yet")
	}

	r.evChan = make(chan string)
	r.abort = make(chan struct{})

	go r.wrld.Render(r.evChan, r.abort, renderDir)

	return r.evChan, nil
}

func (r *Renderer) LoadWorldFromFile(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %s", filename, err)
	}
	defer fh.Close()

	return r.LoadWorldFromReader(fh)
}

func (r *Renderer) LoadWorldFromReader(rd io.Reader) error {
	r.wrld = new(World)
	err := json.NewDecoder(rd).Decode(&r.wrld)
	if err != nil {
		return fmt.Errorf("failed to decode world: %s", err)
	}
	return nil
}
