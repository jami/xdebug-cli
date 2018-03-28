

all: build

build:
	go build -o ./bin/xdebug-cli main.go

@PHONY: build