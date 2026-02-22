package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.5.2"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "BanForge version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BanForge version:", version)
	},
}
