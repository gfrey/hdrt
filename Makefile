.PHONY: default build dev render

default: build

build:
	go get ./...
	go test ./...
	go vet ./...

dev: build
	go-reload hdrt server
	

render: build
	hdrt render default_scene.json