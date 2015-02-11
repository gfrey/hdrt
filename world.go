package hdrt

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"path"
	"sync"
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

	if wrld.Viewplane.Width == 0 || wrld.Viewplane.Height == 0 {
		return fmt.Errorf("width or height (%dx%d)", wrld.Viewplane.Width, wrld.Viewplane.Height)
	}

	return nil
}

const NUM_PARALLEL = 8

func (wrld *World) Render(evChan chan<- string, abortChan <-chan struct{}, renderDir string) {
	pixelInChan := make(chan *pixel)
	pixelOutChan := make(chan *pixel)
	go func(pc chan<- *pixel, ac <-chan struct{}) {
		for y := 0; y < wrld.Viewplane.Height; y++ {
			for x := 0; x < wrld.Viewplane.Width; x++ {
				select {
				case <-ac:
					log.Printf("aborting pixel generator")
					close(pc)
					return
				default:
					pc <- &pixel{x: x, y: y, dir: wrld.dirForPixel(x, y)}
				}
			}
		}
		close(pc)
	}(pixelInChan, abortChan)

	go func() {
		wg := new(sync.WaitGroup)
		for i := 0; i < NUM_PARALLEL; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, pinc <-chan *pixel, poutc chan<- *pixel, ac <-chan struct{}, i int) {
				log.Printf("starting worker %d", i)
				defer wg.Done()
				for pxl := range pinc {
					select {
					case <-ac:
						log.Printf("aborting pixel worker %d", i)
						return
					default:
						wrld.renderPixel(pxl)
						poutc <- pxl
					}
				}
			}(wg, pixelInChan, pixelOutChan, abortChan, i)
		}
		log.Printf("waiting for wait group")
		wg.Wait()
		log.Printf("wait group closed")
		close(pixelOutChan)
	}()

	img := image.NewRGBA(image.Rect(0, 0, wrld.Viewplane.Width, wrld.Viewplane.Height))
	ticker := time.NewTicker(2 * time.Second)
RENDER_LOOP:
	for {
		select {
		case pxl, ok := <-pixelOutChan:
			if !ok {
				break RENDER_LOOP
			}
			img.Set(pxl.x, pxl.y, pxl.col)
		case <-ticker.C:
			filename, err := imgSave(renderDir, img)
			switch err {
			case nil:
				evChan <- path.Base(filename)
			default:
				log.Printf("ERR: %s", err)
			}
		}
	}
	filename, err := imgSave(renderDir, img)
	if err != nil {
		log.Printf("ERR: %s", err)
	}
	evChan <- path.Base(filename)
	close(evChan)
}

func imgSave(renderDir string, img *image.RGBA) (string, error) {
	fh, err := ioutil.TempFile(renderDir, "img")
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}

	err = png.Encode(fh, img)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}

	filename := fh.Name()
	err = fh.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file: %s", err)
	}

	return filename, nil

}

func (wrld *World) renderPixel(pxl *pixel) {
	time.Sleep(50 * time.Millisecond)
	pxl.col = &color.RGBA{200, 0, 0, 255}
}

func (wrld *World) dirForPixel(x, y int) *Vector {
	return nil
}
