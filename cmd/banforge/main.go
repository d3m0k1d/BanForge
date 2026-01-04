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

func Init() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
