package command

import (
	"os"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"github.com/spf13/cobra"
)

var BanListCmd = &cobra.Command{
	Use:   "list",
	Short: "List banned IP adresses",
	Run: func(cmd *cobra.Command, args []string) {
		var log = logger.New(false)
		d, err := storage.NewBanReader()
		if err != nil {
			log.Error("Failed to create database", "error", err)
			os.Exit(1)
		}
		err = d.BanList()
		if err != nil {
			log.Error("Failed to get ban list", "error", err)
			os.Exit(1)
		}
	},
}
