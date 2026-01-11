package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		err := os.Mkdir("/var/log/banforge", 0750)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.Mkdir("/etc/banforge", 0750)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
