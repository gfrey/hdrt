package hdrt

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gfrey/hdrt/vec"
)

type World struct {
	Camera    *Camera
	Viewplane *Viewplane
	Scene     *Scene

	abort chan struct{}
}

func (wrld *World) render(img *image.RGBA) {
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			select {
			case <-wrld.abort:
				return
			default:
				// This is where the magic happens: send ray to scene and determine output color.
				pos, dir := wrld.posAndDirForPixel(x, y)
				img.SetRGBA(x, y, *(wrld.Scene.Render(pos, dir)))
			}
		}
	}
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

func (wrld *World) Render() <-chan *image.RGBA {
	wrld.abort = make(chan struct{})

	imgInChan := make(chan *image.RGBA)
	imgOutChan := make(chan *image.RGBA)
	go func(pc chan<- *image.RGBA, ac <-chan struct{}) {
		var (
			xstride = wrld.Viewplane.ResX / (2 * NUM_PARALLEL)
			ystride = wrld.Viewplane.ResY / (2 * NUM_PARALLEL)
		)

		for y := 0; y < wrld.Viewplane.ResY; y += ystride {
			for x := 0; x < wrld.Viewplane.ResX; x += xstride {

				xmax, ymax := x+xstride, y+ystride
				if xmax > wrld.Viewplane.ResX {
					xmax = wrld.Viewplane.ResX
				}
				if ymax > wrld.Viewplane.ResY {
					ymax = wrld.Viewplane.ResY
				}
				si := image.NewRGBA(image.Rect(x, y, xmax, ymax))

				select {
				case <-ac:
					log.Printf("aborting img generator")
					close(pc)
					return
				default:
					pc <- si
				}
			}
		}
		close(pc)
	}(imgInChan, wrld.abort)

	go func() {
		wg := new(sync.WaitGroup)
		for i := 0; i < NUM_PARALLEL; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, pinc <-chan *image.RGBA, poutc chan<- *image.RGBA, ac <-chan struct{}, i int) {
				log.Printf("starting worker %d", i)
				defer wg.Done()
				for si := range pinc {
					wrld.render(si)
					poutc <- si
				}
			}(wg, imgInChan, imgOutChan, wrld.abort, i)
		}
		log.Printf("waiting for wait group")
		wg.Wait()
		log.Printf("wait group closed")
		close(imgOutChan)
	}()
	return imgOutChan
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