package hdrt

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
)

type World struct {
	Camera    *Camera
	Viewplane *Viewplane
	Scene     *Scene
}

func (wrld *World) Render(evChan chan<- string, abortChan <-chan struct{}, path string) {
	for i := 0; i < 10; i++ {
		select {
		case <-abortChan:
			log.Printf("aborting computation")
			return
		default:
			err := func() error {
				fh, err := ioutil.TempFile(path, "img")
				if err != nil {
					return fmt.Errorf("failed to create file: %s", err)
				}

				img := image.NewRGBA(image.Rect(0, 0, 400, 200))
				for x := 0; x < 400; x++ {
					for y := 0; y < 200; y++ {
						img.Set(x, y, color.RGBA{uint8(255 / i), uint8(255 / i), uint8(255 / i), 255})
					}
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

				evChan <- filename

				return nil
			}()
			if err != nil {
				log.Printf("ERR: %s", err)
			}
		}
	}
	close(evChan)
}
