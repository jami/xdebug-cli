package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jami/xdebug-cli/cfg"
	"github.com/jami/xdebug-cli/dbgp"
	"github.com/jami/xdebug-cli/view"
	"github.com/spf13/cobra"
)

// runCmd represents the hello command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "starts xdbg with script",
	Long:  `starts xdbg with script`,
	Run: func(cmd *cobra.Command, args []string) {
		runRunCmd(args)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func runRunCmd(args []string) {
	server := dbgp.NewServer(CLIArgs.Host, CLIArgs.Port)
	if err := server.Listen(); err != nil {
		fmt.Printf("Error while starting server. %s\n", err)
		os.Exit(1)
	}

	view := view.NewView()
	view.PrintApplicationInformation(cfg.Version, CLIArgs.Host, CLIArgs.Port)

	go startScriptEngine(args)
	server.Accept(listenAccept)
}

func startScriptEngine(args []string) error {
	execCmd := args[0]
	execArgs := []string{}
	// inject connection parameter
	execArgs = append(execArgs, fmt.Sprintf("-dxdebug.remote_host=%s", CLIArgs.Host))
	execArgs = append(execArgs, fmt.Sprintf("-dxdebug.remote_port=%d", CLIArgs.Port))
	// inject rest of args
	execArgs = append(execArgs, args[1:]...)
	// prepare command
	cmd := exec.Command(execCmd, execArgs...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	// exec
	err := cmd.Start()
	if err != nil {
		return err
	}
	// wait until process is done
	return cmd.Wait()
}
