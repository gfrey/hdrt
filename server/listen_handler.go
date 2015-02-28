package server

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gfrey/hdrt"
	"golang.org/x/net/websocket"
)

func ListenHandler(renderDir string) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		err := func() error {
			logger.Printf("web socket connection!! renderDir: %s ", renderDir)

			var msg string
			var wrld *hdrt.World
			for {
				websocket.Message.Receive(ws, &msg)
				switch {
				case strings.HasPrefix(msg, "ABORT"):
					if wrld != nil {
						wrld.Abort()
					}
				case strings.HasPrefix(msg, "CFG"):
					if wrld != nil {
						wrld.Abort()
					}
					go renderImg(ws, msg, renderDir, &wrld)
				}
			}

			return nil
		}()

		if err != nil {
			WSError(ws, err)
		}
	})
}

func renderImg(ws *websocket.Conn, msg, renderDir string, wrldPtr **hdrt.World) {
	var err error
	buf := strings.NewReader(msg[3:])
	*wrldPtr, err = hdrt.LoadWorldFromReader(buf)
	if err != nil {
		WSError(ws, err)
		return
	}

	evChan := make(chan string)
	go func(wrld *hdrt.World, evChan chan<- string) {
		sImgOutChan := wrld.Render()
		img := image.NewRGBA(image.Rect(0, 0, wrld.Viewplane.ResX, wrld.Viewplane.ResY))
		ticker := time.NewTicker(2 * time.Second)
	RENDER_LOOP:
		for {
			select {
			case si, ok := <-sImgOutChan:
				if !ok {
					break RENDER_LOOP
				}
				i := 0
				for y := si.Rect.Min.Y; y < si.Rect.Max.Y; y++ {
					for x := si.Rect.Min.X; x < si.Rect.Max.X; x++ {
						img.SetRGBA(x, y, color.RGBA{si.Buf[i+0], si.Buf[i+1], si.Buf[i+2], si.Buf[i+3]})
						i += 4
					}
				}
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
	}(*wrldPtr, evChan)

	for imagePath := range evChan {
		logger.Printf("send image path %s", imagePath)
		websocket.Message.Send(ws, "IMG"+imagePath)
	}
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