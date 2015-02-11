package server

import (
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
			r := new(hdrt.Renderer)
			for {
				websocket.Message.Receive(ws, &msg)
				switch {
				case strings.HasPrefix(msg, "ABORT"):
					r.Abort()
				case strings.HasPrefix(msg, "CFG"):
					r.Abort()
					go func(msg string) {
						buf := strings.NewReader(msg[3:])
						err := r.LoadWorldFromReader(buf)
						if err != nil {
							WSError(ws, err)
							return
						}

						logger.Printf("--> start rendering")
						ch, err := r.Render(renderDir)
						if err != nil {
							WSError(ws, err)
							return
						}

						for imagePath := range ch {
							logger.Printf("send image path %s", imagePath)
							websocket.Message.Send(ws, "IMG"+imagePath)
						}
					}(msg)
				}
			}

			return nil
		}()

		if err != nil {
			WSError(ws, err)
		}
	})
}
