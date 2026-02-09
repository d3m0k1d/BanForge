package main

import (
	"fmt"
	"os"

	"github.com/d3m0k1d/BanForge/cmd/banforge/command"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "banforge",
	Short: "IPS log-based written on Golang",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Init() {

}

func Execute() {
	rootCmd.AddCommand(command.DaemonCmd)
	rootCmd.AddCommand(command.InitCmd)
	rootCmd.AddCommand(command.RuleCmd)
	rootCmd.AddCommand(command.BanCmd)
	rootCmd.AddCommand(command.UnbanCmd)
	rootCmd.AddCommand(command.BanListCmd)
	rootCmd.AddCommand(command.VersionCmd)
	command.RuleRegister()
	command.FwRegister()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
