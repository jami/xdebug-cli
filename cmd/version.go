package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the hello command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "shows the version of the application",
	Long:  `shows the version of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version", CLIArgs.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
