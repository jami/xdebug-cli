

all: build

build: test
	#go build -o ./bin/xdebug-cli main.go
	go build -o ./bin/xdebug-cli cmd/xdebug-cli.go

test:
	go test -v ./...


@PHONY: build test