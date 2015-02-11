package server

import (
	"bytes"
	"net/http"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(RootTPL))

	if err != nil {
		HTTPError(w, err, 500)
	}

	return
}

const RootTPL = `
<html>
	<head>
		<title>Ray Tracer</title>
		<!-- Latest compiled and minified CSS -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">

		<!-- Optional theme -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">

		<!-- Latest compiled and minified JavaScript -->
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
	</head>
	<body>
		<h1>Hallo Welt</h1>
		<p>Hier kommt der Keks</p>
	</body>
</html>
`

func render(t string, i interface{}) (string, error) {
	tpl, err := template.New(t).Parse(t)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, i)
	return buf.String(), err
}
