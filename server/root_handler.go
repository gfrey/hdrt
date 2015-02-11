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

		<script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>

		<!-- Latest compiled and minified CSS -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">

		<!-- Optional theme -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">

		<!-- Latest compiled and minified JavaScript -->
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>

		<link rel="stylesheet" href="/styles.css">
		<script src="/app.js"></script>
	</head>
	<body>
		<div class="container">
			<div class="page-header">
				<h1>Ray Tracing</h1>
				<p class="lead">Hier kommt der Keks</p>
			</div>
			<div class="row">
				<div class="col-md-6">
					<form>
						<textarea id="editor"></textarea>
						<button id="render-btn" class="btn btn-primary">Render</button>
					</form>
				</div>
				<div class="col-md-6">
					<div id="scene">
						scene
					</div>
				</div>
			</div>
		</div>
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
