package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the hello command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts xdbg in listening mode",
	Long:  `starts xdbg in listening mode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
