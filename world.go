package hdrt

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"path"
	"time"
)

type World struct {
	Camera    *Camera
	Viewplane *Viewplane
	Scene     *Scene
}

type pixel struct {
	x, y int
	col  *color.RGBA
	dir  *Vector
}

func (wrld *World) Validate() error {
	if wrld.Viewplane == nil {
		return fmt.Errorf("viewplane not set")
	}
	return nil
}

const NUM_PARALLEL = 8

func (wrld *World) Render(evChan chan<- string, abortChan <-chan struct{}, renderDir string) {
	pixelInChan := make(chan *pixel)
	pixelOutChan := make(chan *pixel)
	go func(pc chan<- *pixel) {
		for x := 0; x < wrld.Viewplane.Width; x++ {
			for y := 0; y < wrld.Viewplane.Height; y++ {
				pc <- &pixel{x: x, y: y, dir: wrld.dirForPixel(x, y)}
			}
		}
	}(pixelInChan)

	for i := 0; i < NUM_PARALLEL; i++ {
		go func(pinc <-chan *pixel, poutc chan<- *pixel) {
			for pxl := range pinc {
				wrld.renderPixel(pxl)
				poutc <- pxl
			}
		}(pixelInChan, pixelOutChan)
	}

	img := image.NewRGBA(image.Rect(0, 0, wrld.Viewplane.Width, wrld.Viewplane.Height))
	ticker := time.Tick(2 * time.Second)
RENDER_LOOP:
	for {
		select {
		case <-abortChan:
			close(pixelInChan)
			log.Printf("aborting computation")
			return
		case pxl, ok := <-pixelOutChan:
			if !ok {
				break RENDER_LOOP
			}
			img.Set(pxl.x, pxl.y, pxl.col)
		case <-ticker:
			err := func() error {
				fh, err := ioutil.TempFile(renderDir, "img")
				if err != nil {
					return fmt.Errorf("failed to create file: %s", err)
				}

				err = png.Encode(fh, img)
				if err != nil {
					return fmt.Errorf("failed to create file: %s", err)
				}

				filename := fh.Name()
				err = fh.Close()
				if err != nil {
					return fmt.Errorf("failed to close file: %s", err)
				}

				evChan <- path.Base(filename)
				return nil
			}()
			if err != nil {
				log.Printf("ERR: %s", err)
			}
		}
	}
	close(evChan)
}

func (wrld *World) renderPixel(pxl *pixel) {
	pxl.col = &color.RGBA{200, 0, 0, 255}
}

func (wrld *World) dirForPixel(x, y int) *Vector {
	return nil
}
