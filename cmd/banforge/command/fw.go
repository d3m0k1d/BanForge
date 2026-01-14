package command

import (
	"fmt"
	"net"
	"os"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/spf13/cobra"
)

var (
	ip string
)
var UnbanCmd = &cobra.Command{
	Use:   "unban",
	Short: "Unban IP",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fw := cfg.Firewall.Name
		b := blocker.GetBlocker(fw, cfg.Firewall.Config)
		if ip == "" {
			fmt.Println("IP can't be empty")
			os.Exit(1)
		}
		if net.ParseIP(ip) == nil {
			fmt.Println("Invalid IP")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = b.Unban(ip)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("IP unblocked successfully!")
	},
}

var BanCmd = &cobra.Command{
	Use:   "ban",
	Short: "Ban IP",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fw := cfg.Firewall.Name
		b := blocker.GetBlocker(fw, cfg.Firewall.Config)
		if ip == "" {
			fmt.Println("IP can't be empty")
			os.Exit(1)
		}
		if net.ParseIP(ip) == nil {
			fmt.Println("Invalid IP")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = b.Ban(ip)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("IP unblocked successfully!")
	},
}

func FwRegister() {
	BanCmd.Flags().StringVarP(&ip, "ip", "i", "", "ip to ban")
	UnbanCmd.Flags().StringVarP(&ip, "ip", "i", "", "ip to unban")
}
