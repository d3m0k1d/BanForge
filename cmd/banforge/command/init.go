package command

import (
	"fmt"
	"os"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize BanForge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing BanForge...")

		err := config.CreateConf()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = config.FindFirewall()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b := blocker.GetBlocker(cfg.Firewall.Name, cfg.Firewall.Config)
		err = b.Setup(cfg.Firewall.Config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Firewall configured")

		err = storage.CreateTables()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Firewall detected and configured")

		fmt.Println("BanForge initialized successfully!")
	},
}
