package main

import (
	"math/rand"
	"time"

	"github.com/jami/xdebug-cli/cmd"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cmd.Execute()
}
