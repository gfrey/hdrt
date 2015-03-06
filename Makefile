.PHONY: default build dev render

default: build

build:
	go get ./...
	go test ./...
	go vet ./...

dev: build
	go-reload hdrt server
	

render:
	go get ./...
	hdrt render default_scene.json