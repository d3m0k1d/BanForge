package blocker

import (
	"fmt"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type BlockerEngine interface {
	Ban(ip string) error
	Unban(ip string) error
	Setup(config string) error
	PortOpen(port int, protocol string) error
	PortClose(port int, protocol string) error
}

func GetBlocker(fw string, config string) (BlockerEngine, error) {
	switch fw {
	case "ufw":
		return NewUfw(logger.New(false)), nil
	case "iptables":
		return NewIptables(logger.New(false), config), nil
	case "nftables":
		return NewNftables(logger.New(false), config), nil
	case "firewalld":
		return NewFirewalld(logger.New(false)), nil
	default:
		return nil, fmt.Errorf("unknown firewall: %s", fw)
	}
}
