package hdrt

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gfrey/hdrt/vec"
)

type World struct {
	Camera    *Camera
	Viewplane *Viewplane
	Scene     *Scene

	evChan chan string
	abort  chan struct{}
}

type pixel struct {
	x, y     int
	col      *color.RGBA
	pos, dir *vec.Vector
}

func LoadWorldFromFile(filename string) (*World, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %s", filename, err)
	}
	defer fh.Close()

	return LoadWorldFromReader(fh)
}

func LoadWorldFromReader(rd io.Reader) (*World, error) {
	wrld := new(World)
	err := json.NewDecoder(rd).Decode(&wrld)
	if err != nil {
		return nil, fmt.Errorf("failed to decode world: %s", err)
	}
	return wrld, wrld.Init()
}

func (wrld *World) Init() error {
	err := wrld.validate()
	if err != nil {
		return err
	}

	err = wrld.Camera.Init()
	if err != nil {
		return err
	}
	err = wrld.Viewplane.Init(wrld.Camera)
	if err != nil {
		return err
	}
	return nil
}

func (wrld *World) validate() error {
	if wrld.Viewplane == nil {
		return fmt.Errorf("viewplane not set")
	}

	if wrld.Viewplane.ResX == 0 || wrld.Viewplane.ResY == 0 {
		return fmt.Errorf("width or height (%dx%d)", wrld.Viewplane.ResX, wrld.Viewplane.ResY)
	}

	if wrld.Camera == nil {
		return fmt.Errorf("camera not set")
	}

	if wrld.Camera.Position == nil {
		return fmt.Errorf("camera position not set")
	}

	if wrld.Scene == nil {
		return fmt.Errorf("scene not set")
	}

	if len(wrld.Scene.Objects) == 0 {
		return fmt.Errorf("no objects in scene")
	}

	return nil
}

const NUM_PARALLEL = 8

func (wrld *World) Render() <-chan *pixel {
	wrld.evChan = make(chan string)
	wrld.abort = make(chan struct{})

	pixelInChan := make(chan *pixel)
	pixelOutChan := make(chan *pixel)
	go func(pc chan<- *pixel, ac <-chan struct{}) {
		for y := 0; y < wrld.Viewplane.ResY; y++ {
			for x := 0; x < wrld.Viewplane.ResX; x++ {
				select {
				case <-ac:
					log.Printf("aborting pixel generator")
					close(pc)
					return
				default:
					pos, dir := wrld.posAndDirForPixel(x, y)
					pc <- &pixel{x: x, y: y, pos: pos, dir: dir}
				}
			}
		}
		close(pc)
	}(pixelInChan, wrld.abort)

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
						// This is where the magic happens: send ray to scene and determine output color.
						pxl.col = wrld.Scene.Render(pxl.pos, pxl.dir)
						poutc <- pxl
					}
				}
			}(wg, pixelInChan, pixelOutChan, wrld.abort, i)
		}
		log.Printf("waiting for wait group")
		wg.Wait()
		log.Printf("wait group closed")
		close(pixelOutChan)
	}()
	return pixelOutChan
}

func (wrld *World) RenderToWeb(renderDir string) <-chan string {

	pixelOutChan := wrld.Render()
	go func() {
		img := image.NewRGBA(image.Rect(0, 0, wrld.Viewplane.ResX, wrld.Viewplane.ResY))
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
					wrld.evChan <- path.Base(filename)
				default:
					log.Printf("ERR: %s", err)
				}
			}
		}
		filename, err := imgSave(renderDir, img)
		if err != nil {
			log.Printf("ERR: %s", err)
		}
		wrld.evChan <- path.Base(filename)
		close(wrld.evChan)
	}()
	return wrld.evChan
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

func (wrld *World) posAndDirForPixel(x, y int) (*vec.Vector, *vec.Vector) {
	positionPixel := wrld.Viewplane.PixelPosition(x, y)
	dir := vec.Add(positionPixel, vec.ScalarMultiply(wrld.Camera.Position, -1.0))
	dir.Normalize()
	return positionPixel, dir
}

func (wrld *World) Abort() {
	if wrld.abort != nil {
		log.Printf("closing abort channel")
		close(wrld.abort)
		wrld.abort = nil
	}
}