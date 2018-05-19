PACKAGES=$(shell go list ./... | grep -v /vendor/)
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))

BINARY=xdbg

# target
VERSION ?= `git rev-parse --short HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=github.com/jami/xdebug-cli/cmd.version=$(VERSION)"

all: build/local

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make deps     		Install go deps.'
	@echo '    make build    		Compile the project.'
	@echo '    make build/docker	Restore all dependencies.'
	@echo '    make restore  		Restore all dependencies.'
	@echo '    make clean    		Clean the directory tree.'
	@echo

test: ## run tests, except integration tests
	@go test ${RACE} ${PACKAGES}

deps:
	go get -u github.com/tcnksm/ghr
	go get -u github.com/mitchellh/gox
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/mitchellh/go-homedir
	go get -u github.com/spf13/cobra
	go get -u github.com/spf13/viper

build:
	@echo "Compiling all arch targets ..."
	@mkdir -p ./bin
	@gox $(LDFLAGS) -output "bin/${BINARY}_{{.OS}}_{{.Arch}}" -os="linux" -os="darwin" -arch="386" -arch="amd64" ./
	@echo "All done! The binaries is in ./bin let's have fun!"

build/local:
	@echo "Compiling local arch ..."
	@mkdir -p ./bin
	@go build $(LDFLAGS) -i -o "./bin/${BINARY}"

build/release:
	@echo "Compiling..."
	@mkdir -p ./bin
	@gox $(LDFLAGS) -tags netgo -ldflags '-w -extldflags "-static"' -output "bin/${BINARY}_{{.OS}}_{{.Arch}}" -os="linux" -os="darwin" -arch="386" -arch="amd64" ./
	@echo "All done! The binaries is in ./bin let's have fun!"
 
build/docker: build
	@docker build -t xdbg-example:latest .

vet: ## run go vet
	@test -z "$$(go vet ${PACKAGES} 2>&1 | grep -v '*composite literal uses unkeyed fields|exit status 0)' | tee /dev/stderr)"

ci: vet test

restore:
	@dep ensure