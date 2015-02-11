package server

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/websocket"
)

func Router(renderDir string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", IndexHandler)
	mux.Handle("/listen", ListenHandler(renderDir))
	registerStaticHandlers(mux, "assets")

	mux.Handle("/renders/", http.StripPrefix("/renders", http.FileServer(http.Dir(renderDir))))

	return mux
}

func HTTPError(w http.ResponseWriter, err error, status int) {
	logger.Printf("ERROR status=%d %s err=%s", status, http.StatusText(status), err.Error())
	w.WriteHeader(status)
	body := fmt.Sprintf("status=%d %s err=%s", status, http.StatusText(status), err)
	w.Write([]byte(body))
}

func WSError(ws *websocket.Conn, err error) {
	logger.Printf("ERROR on websocket: %s", err.Error())
	websocket.Message.Send(ws, "ERR"+err.Error())
}

func registerStaticHandlers(mux *http.ServeMux, folder string) {
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		route := strings.TrimPrefix(path, folder)
		mux.HandleFunc(route, staticFileHander(path))
		return nil
	})
}

func staticFileHander(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(path)
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))

		http.ServeFile(w, r, path)
	}
}
