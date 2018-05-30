package cfg

// Version global
var Version string = "develop"

// CLIParameter holds the global flagset
type CLIParameter struct {
	Host    string
	Port    uint16
	Trigger string
	Version string
}
