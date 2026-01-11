package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "banforge",
	Short: "IPS log-based written on Golang",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize BanForge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing BanForge...")
		os.Mkdir("/var/log/banforge", 0755)
		os.Mkdir("/etc/banforge", 0755)
	},
}

func Init() {

}

func Execute() {
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
