package server

import (
	"net/http"

	"github.com/gfrey/hdrt"
	"golang.org/x/net/websocket"
)

func ListenHandler(renderDir string) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		err := func() error {
			logger.Printf("web socket connection!! renderDir: %s ", renderDir)

			r := new(hdrt.Renderer)
			err := r.LoadWorldFromReader(ws)
			if err != nil {
				return err
			}

			ch, err := r.Render(renderDir)
			if err != nil {
				return err
			}

			for imagePath := range ch {
				logger.Printf("got image path %s", imagePath)
				ws.Write([]byte(imagePath))
			}

			return nil
		}()

		if err != nil {
			WSError(ws, err)
		}
	})
}
