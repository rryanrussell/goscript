.PHONY: build install test example dommer-libdom

build:
	go build -o bin/goscript ./cmd/goscript
	go build -o bin/dommer ./cmd/dommer

install:
	go install ./cmd/goscript
	go install ./cmd/dommer

test:
	go test ./pkg/... ./cmd/...

example: build
	./bin/goscript examples/webgl-app

lib:
	mkdir -p lib

lib/lib.dom.d.ts: lib
	curl -L -o lib/lib.dom.d.ts https://unpkg.com/typescript@latest/lib/lib.dom.d.ts

dommer-libdom: build lib/lib.dom.d.ts
	./bin/dommer -output lib/dom.extern.go
