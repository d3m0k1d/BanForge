package config

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func CreateConf() {
	if syscall.Geteuid() != 0 {
		os.Exit(1)
		fmt.Printf("You must be root to run\n, use the sudo/doas")
	}
	exec.Command("mkdir /etc/banforge")
	exec.Command("touch /etc/banforge/config.toml")

}

func CheckSysConf() {
}
