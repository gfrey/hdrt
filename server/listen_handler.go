package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gfrey/hdrt"
	"golang.org/x/net/websocket"
)

func ListenHandler(renderDir string) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		err := func() error {
			logger.Printf("web socket connection!! renderDir: %s ", renderDir)

			var msg string
			var r *hdrt.Renderer
			for {
				websocket.Message.Receive(ws, &msg)
				switch {
				case strings.HasPrefix(msg, "ABORT"):
					log.Printf("aborting everything")
					if r != nil {
						r.Abort()
					}
					return nil
				case strings.HasPrefix(msg, "CFG"):
					r = new(hdrt.Renderer)

					buf := strings.NewReader(msg[3:])
					err := r.LoadWorldFromReader(buf)
					if err != nil {
						WSError(ws, err)
						continue
					}

					logger.Printf("--> start rendering")
					ch, err := r.Render(renderDir)
					if err != nil {
						WSError(ws, err)
						continue
					}

					for imagePath := range ch {
						logger.Printf("send image path %s", imagePath)
						websocket.Message.Send(ws, "IMG"+imagePath)
					}
				}
			}

			return nil
		}()

		if err != nil {
			WSError(ws, err)
		}
	})
}
