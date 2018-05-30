package cmd

import (
	"fmt"
	"os"

	"github.com/jami/xdebug-cli/cfg"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DefaultConfigFile name
const DefaultConfigFile = ".xdbgcli"

var (
	// RootCmd main command
	RootCmd = &cobra.Command{
		Use:   "xdbg",
		Short: "xdebug cli",
		Long:  `CLI debugger for php applications`,
	}
	// cfgFile
	cfgFile string
	// CLIArgs global parameter
	CLIArgs = &cfg.CLIParameter{}

	version string
)

func init() {
	CLIArgs.Version = version

	cobra.OnInitialize(initConfigFile)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xdbgcli.yaml)")
	RootCmd.PersistentFlags().StringVarP(&CLIArgs.Host, "host", "l", "127.0.0.1", "Listener host")
	RootCmd.PersistentFlags().Uint16VarP(&CLIArgs.Port, "port", "p", 9000, "Listener port")
	RootCmd.PersistentFlags().StringVarP(&CLIArgs.Trigger, "trigger", "t", "", "Debug session trigger")
}

func initConfigFile() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(DefaultConfigFile)
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute cli arg process
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
