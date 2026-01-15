package blocker

import (
	"fmt"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type BlockerEngine interface {
	Ban(ip string) error
	Unban(ip string) error
	Setup(config string) error
}

func GetBlocker(fw string, config string) BlockerEngine {
	switch fw {
	case "ufw":
		return NewUfw(logger.New(false))
	case "iptables":
		return NewIptables(logger.New(false), config)
	case "nftables":
		return NewNftables(logger.New(false), config)
	case "firewalld":
		return NewFirewalld(logger.New(false))
	default:
		panic(fmt.Sprintf("Unknown firewall: %s", fw))
	}
}
