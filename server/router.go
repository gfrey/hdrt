package server

import "net/http"

func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", RootHandler)

	return mux
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Miau"))

	if err != nil {
		HTTPError(w, err, 500)
	}

	return
}

func HTTPError(w http.ResponseWriter, err error, status int) {
	logger.Printf("ERROR status=%d %s err=%s", status, http.StatusText(status), err.Error())
	w.WriteHeader(status)
	w.Write([]byte("error yeah"))
}
