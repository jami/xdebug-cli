package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the hello command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "starts xdbg with script",
	Long:  `starts xdbg with script`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
