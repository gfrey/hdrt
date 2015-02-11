default: build

build:
	go get ./...

dev: build
	go-reload hdrt server
