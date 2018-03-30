package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jami/xdebug-cli"
)

var (
	listenAddr string
	listenPort int
)

func parseFlags() error {
	flag.StringVar(&listenAddr, "host", "127.0.0.1", "specify host to listen to")
	flag.IntVar(&listenPort, "port", 9000, "specify port to listen to")

	flag.Parse()

	return nil
}

func main() {
	if err := parseFlags(); err != nil {
		fmt.Printf("Error while parsing cli arguments. %s\n", err)
		os.Exit(1)
	}

	server := xdebugcli.NewDBGPServer(listenAddr, listenPort)
	if err := server.Listen(); err != nil {
		fmt.Printf("Error while starting server. %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listening to %s:%d\n", listenAddr, listenPort)
	server.Accept()
}
