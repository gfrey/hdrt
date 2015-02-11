package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"
)

var RootTPL = func() []byte {
	tpl, err := ioutil.ReadFile("tpl/root.html")
	if err != nil {
		panic("could not read tpl/root.html")
	}
	return tpl
}()

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(RootTPL)

	if err != nil {
		HTTPError(w, err, 500)
	}

	return
}

func render(t string, i interface{}) (string, error) {
	tpl, err := template.New(t).Parse(t)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, i)
	return buf.String(), err
}
