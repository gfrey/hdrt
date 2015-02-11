package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"
)

var IndexTPL = func() []byte {
	tpl, err := ioutil.ReadFile("tpl/index.html")
	if err != nil {
		panic("could not read tpl/index.html")
	}
	return tpl
}()

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(IndexTPL)

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
