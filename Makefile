

all: build

build:
	#go build -o ./bin/xdebug-cli main.go
	go build -o ./bin/xdebug-cli cmd/xdebug-cli.go

@PHONY: build