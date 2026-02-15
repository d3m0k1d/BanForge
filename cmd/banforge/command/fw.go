package command

import (
	"fmt"
	"net"
	"os"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"github.com/spf13/cobra"
)

var (
	ttl_fw   string
	port     int
	protocol string
)

var UnbanCmd = &cobra.Command{
	Use:   "unban",
	Short: "Unban IP",
	Run: func(cmd *cobra.Command, args []string) {
		err := func() error {
			if len(args) == 0 {
				return fmt.Errorf("IP can't be empty")
			}
			if ttl_fw == "" {
				ttl_fw = "1y"
			}
			ip := args[0]
			db, err := storage.NewBanWriter()
			if err != nil {
				return err
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			fw := cfg.Firewall.Name
			b := blocker.GetBlocker(fw, cfg.Firewall.Config)
			if ip == "" {
				return fmt.Errorf("IP can't be empty")
			}
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("invalid IP")
			}
			if err != nil {
				return err
			}
			err = b.Unban(ip)
			if err != nil {
				return err
			}
			err = db.RemoveBan(ip)
			if err != nil {
				return err
			}
			fmt.Println("IP unblocked successfully!")
			return nil
		}()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var BanCmd = &cobra.Command{
	Use:   "ban",
	Short: "Ban IP",
	Run: func(cmd *cobra.Command, args []string) {
		err := func() error {
			if len(args) == 0 {
				return fmt.Errorf("IP can't be empty")
			}
			if ttl_fw == "" {
				ttl_fw = "1y"
			}
			ip := args[0]
			db, err := storage.NewBanWriter()
			if err != nil {
				return err
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			fw := cfg.Firewall.Name
			b := blocker.GetBlocker(fw, cfg.Firewall.Config)
			if ip == "" {
				return fmt.Errorf("IP can't be empty")
			}
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("invalid IP")
			}
			if err != nil {
				return err
			}
			err = b.Ban(ip)
			if err != nil {
				return err
			}
			err = db.AddBan(ip, ttl_fw, "manual ban")
			if err != nil {
				return err
			}
			fmt.Println("IP blocked successfully!")
			return nil
		}()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var PortCmd = &cobra.Command{
	Use:   "port",
	Short: "Ports commands",
}

var PortOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open ports on firewall",
	Run: func(cmd *cobra.Command, args []string) {
		if protocol == "" {
			fmt.Println("Protocol can't be empty")
			os.Exit(1)
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fw := cfg.Firewall.Name
		b := blocker.GetBlocker(fw, cfg.Firewall.Config)
		err = b.PortOpen(port, protocol)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Port opened successfully!")
	},
}

var PortCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close ports on firewall",
	Run: func(cmd *cobra.Command, args []string) {
		if protocol == "" {
			fmt.Println("Protocol can't be empty")
			os.Exit(1)
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fw := cfg.Firewall.Name
		b := blocker.GetBlocker(fw, cfg.Firewall.Config)
		err = b.PortClose(port, protocol)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Port closed successfully!")
	},
}

func FwRegister() {
	BanCmd.Flags().StringVarP(&ttl_fw, "ttl", "t", "", "ban time")
	PortCmd.AddCommand(PortOpenCmd)
	PortCmd.AddCommand(PortCloseCmd)
	PortOpenCmd.Flags().IntVarP(&port, "port", "p", 0, "port number")
	PortOpenCmd.Flags().StringVarP(&protocol, "protocol", "c", "", "protocol")
	PortCloseCmd.Flags().IntVarP(&port, "port", "p", 0, "port number")
	PortCloseCmd.Flags().StringVarP(&protocol, "protocol", "c", "", "protocol")
}
