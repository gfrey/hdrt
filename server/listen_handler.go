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

			var err error
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
					go func(msg string) {
						buf := strings.NewReader(msg[3:])
						wrld, err = hdrt.LoadWorldFromReader(buf)
						if err != nil {
							WSError(ws, err)
							return
						}

						logger.Printf("--> start rendering")
						ch := wrld.RenderToWeb(renderDir)

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
