package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jami/xdebug-cli/cfg"
	"github.com/jami/xdebug-cli/dbgp"
	"github.com/jami/xdebug-cli/view"
	"github.com/spf13/cobra"
)

// listeningCmd represents the dbg listening mode
var listeningCmd = &cobra.Command{
	Use:   "listen",
	Short: "starts xdbg in listening mode",
	Long:  `starts xdbg in listening mode`,
	Run: func(cmd *cobra.Command, args []string) {
		runListeningCmd()
	},
}

var (
	matchBreakLine = regexp.MustCompile("^:(\\d+)(.+)?$")
	matchBreakFile = regexp.MustCompile("^([^:]+):(\\d+)(.+)?$")
)

func init() {
	RootCmd.AddCommand(listeningCmd)
}

func runListeningCmd() {
	server := dbgp.NewServer(CLIArgs.Host, CLIArgs.Port)
	if err := server.Listen(); err != nil {
		fmt.Printf("Error while starting server. %s\n", err)
		os.Exit(1)
	}

	view := view.NewView()
	view.PrintApplicationInformation(cfg.Version, CLIArgs.Host, CLIArgs.Port)

	server.Accept(listenAccept)
}

func listenAccept(c *dbgp.Connection) {
	defer c.Close()

	client := dbgp.NewClient(c)
	if err := client.Init(); err != nil {
		fmt.Println("Connection initialization error: ", err)
		return
	}

	// set per default
	client.SetExceptionBreakpoint()

	view := view.NewView()

	for {
		view.PrintInputPrefix()
		cmdBuffer := strings.TrimSpace(view.GetInputLine())
		cmdSlice := strings.Split(cmdBuffer, " ")

		if len(cmdSlice) == 0 || len(cmdSlice[0]) == 0 {
			view.PrintErrorLn("Unknown command. Try help")
			continue
		}

		switch cmdSlice[0] {
		case "q", "quit":
			view.Print("Terminate debug session? (y or n) ")
			confirmation := strings.TrimSpace(view.GetInputLine())

			if confirmation == "y" {
				view.PrintLn("closing debug session")
				return
			}
		case "h", "help", "?":
			view.ShowHelpMessage()
		case "info", "i":
			getInfo(view, cmdSlice, client)
		case "list", "l":
			showCodeLines(view, client)
		case "break", "b":
			setBreakpoint(view, cmdSlice, client)
		case "run", "r":
			run(view, client)
		case "step", "s":
			step(view, client)
		case "next", "n":
			next(view, client)
		case "finish", "f":
			finish(view, client)
		default:
			view.PrintErrorLn("Unknown command. Try help")
		}

		if client.Session.State == dbgp.StateStopped {
			view.PrintLn("closing debug session")
			return
		}
	}
}

func showCodeLines(view *view.View, c *dbgp.Client) {
	view.PrintSourceLn(c.Session.CurrentFile, c.Session.CurrentLine, 10)
}

func updateState(view *view.View, resp *dbgp.ProtocolResponse, c *dbgp.Client) {
	if resp.Reason != "ok" {
		view.PrintErrorLn("response error reason " + resp.Reason)
		return
	}

	switch resp.Status {
	case "break":
		c.Session.State = dbgp.StateBreak
		c.Session.CurrentFile = resp.Message.Filename
		c.Session.CurrentLine = resp.Message.Line

		if resp.Command == "run" {
			view.PrintLn(
				fmt.Sprintf(
					"breakpoint, %s at %s:%d",
					resp.Message.Exception,
					c.Session.CurrentFile,
					c.Session.CurrentLine))
		}

		view.PrintSourceLn(c.Session.CurrentFile, c.Session.CurrentLine, 1)

		if resp.Message.Exception == "Fatal error" {
			view.PrintErrorLn("fatal error occured")
			view.PrintLn("")
			view.PrintLn(resp.Message.Value)
			view.PrintLn("")
		}
	case "stopping":
		c.Session.State = dbgp.StateStopped
		view.PrintLn("session stopped")
	default:
		view.PrintErrorLn("unknown session state " + resp.Status)
	}
}

func finish(view *view.View, c *dbgp.Client) {
	_, err := c.Finish()
	if err != nil {
		view.PrintErrorLn(err.Error())
		return
	}
	c.Session.State = dbgp.StateStopped
	view.PrintLn("session closed")
}

func step(view *view.View, c *dbgp.Client) {
	resp, err := c.Step()
	if err != nil {
		view.PrintErrorLn(err.Error())
		return
	}

	updateState(view, resp, c)
}

func next(view *view.View, c *dbgp.Client) {
	resp, err := c.Next()
	if err != nil {
		view.PrintErrorLn(err.Error())
		return
	}

	updateState(view, resp, c)
}

func run(view *view.View, c *dbgp.Client) {
	if c.Session.State != dbgp.StateStarting && c.Session.State != dbgp.StateBreak {
		view.PrintErrorLn("debugger is already running")
		return
	}

	view.PrintLn(fmt.Sprintf("starting program: %s\n", c.Session.CurrentFile))
	resp, err := c.Run()
	if err != nil {
		view.PrintErrorLn(err.Error())
		return
	}

	updateState(view, resp, c)
}

func setBreakpoint(view *view.View, args []string, c *dbgp.Client) {
	if len(args) < 2 || args[1] == "help" {
		view.ShowBreakpointHelpMessage()
		return
	}

	t := strings.Join(args[1:], " ")
	file := c.Session.CurrentFile
	line := 0
	expr := ""

	if result := matchBreakFile.FindStringSubmatch(t); result != nil {
		file = strings.TrimSpace(result[1])
		line, _ = strconv.Atoi(result[2])

		if len(result) > 2 {
			expr = strings.TrimSpace(result[3])
		}
	} else if result := matchBreakLine.FindStringSubmatch(t); result != nil {
		line, _ = strconv.Atoi(result[1])

		if len(result) > 1 {
			expr = strings.TrimSpace(result[2])
		}
	}

	if len(expr) > 0 {
		expr = base64.StdEncoding.EncodeToString([]byte(expr))
	}

	if err := c.SetBreakpoint(file, line, expr); err != nil {
		view.PrintErrorLn("error while setting breakpoint. " + err.Error())
	}
}

func getInfo(view *view.View, args []string, c *dbgp.Client) {
	if len(args) < 2 || args[1] == "help" {
		view.ShowInfoHelpMessage()
		return
	}

	switch a := args[1]; a {
	case "breakpoints":
		bpl, err := c.GetBreakpointList()
		if err != nil {
			view.PrintErrorLn(err.Error())
			return
		}

		if len(bpl) == 0 {
			view.PrintLn("No breakpoints available")
			return
		}

		view.ShowInfoBreakpoints(bpl)
	default:
		view.PrintLn("Unknown command: info " + a)
		view.ShowInfoHelpMessage()
	}
}
