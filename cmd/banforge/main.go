package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/judge"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/parser"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"github.com/spf13/cobra"
)

var (
	ip      string
	name    string
	service string
	path    string
	status  string
	method  string
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

		if _, err := os.Stat("/var/log/banforge"); err == nil {
			fmt.Println("/var/log/banforge already exists, skipping...")
		} else if os.IsNotExist(err) {
			err := os.Mkdir("/var/log/banforge", 0750)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Created /var/log/banforge")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
		if _, err := os.Stat("/var/lib/banforge"); err == nil {
			fmt.Println("/var/lib/banforge already exists, skipping...")
		} else if os.IsNotExist(err) {
			err := os.Mkdir("/var/lib/banforge", 0750)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Created /var/lib/banforge")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}

		if _, err := os.Stat("/etc/banforge"); err == nil {
			fmt.Println("/etc/banforge already exists, skipping...")
		} else if os.IsNotExist(err) {
			err := os.Mkdir("/etc/banforge", 0750)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Created /etc/banforge")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}

		err := config.CreateConf()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Config created")

		err = config.FindFirewall()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		db, err := storage.NewDB()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = db.CreateTable()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func() {
			err = db.Close()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
		fmt.Println("Firewall detected and configured")

		fmt.Println("BanForge initialized successfully!")
	},
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run BanForge daemon process",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New(false)
		log.Info("Starting BanForge daemon")
		db, err := storage.NewDB()
		if err != nil {
			log.Error("Failed to create database", "error", err)
			os.Exit(1)
		}
		defer func() {
			err = db.Close()
			if err != nil {
				log.Error("Failed to close database connection", "error", err)
			}
		}()
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
		var b blocker.BlockerEngine
		fw := cfg.Firewall.Name
		b = blocker.GetBlocker(fw, cfg.Firewall.Config)
		r, err := config.LoadRuleConfig()
		if err != nil {
			log.Error("Failed to load rules", "error", err)
			os.Exit(1)
		}
		j := judge.New(db, b)
		j.LoadRules(r)
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				if err := j.ProcessUnviewed(); err != nil {
					log.Error("Failed to process unviewed", "error", err)
				}
			}
		}()

		for _, svc := range cfg.Service {
			log.Info("Processing service", "name", svc.Name, "enabled", svc.Enabled, "path", svc.LogPath)

			if !svc.Enabled {
				log.Info("Service disabled, skipping", "name", svc.Name)
				continue
			}

			if svc.Name != "nginx" {
				log.Info("Only nginx supported, skipping", "name", svc.Name)
				continue
			}

			log.Info("Starting parser for service", "name", svc.Name, "path", svc.LogPath)

			pars, err := parser.NewScanner(svc.LogPath)
			if err != nil {
				log.Error("Failed to create scanner", "service", svc.Name, "error", err)
				continue
			}

			go pars.Start()
			go func(p *parser.Scanner, serviceName string) {
				log.Info("Starting nginx parser", "service", serviceName)
				ng := parser.NewNginxParser()
				resultCh := make(chan *storage.LogEntry, 100)
				ng.Parse(p.Events(), resultCh)
				go storage.Write(db, resultCh)
			}(pars, svc.Name)
		}

		select {}
	},
}

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

// Rule block
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage rules",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "CLI interface for add new rule to file /etc/banforge/rules.toml",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			fmt.Printf("Rule name can't be empty\n")
			os.Exit(1)
		}
		if service == "" && path == "" && status == "" && method == "" {
			fmt.Printf("At least 1 rule field must be filled in.")
			os.Exit(1)
		}

		err := config.NewRule(name, service, path, status, method)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Rule added successfully!")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List rules",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := config.LoadRuleConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, rule := range r {
			fmt.Printf("Name: %s\nService: %s\nPath: %s\nStatus: %s\nMethod: %s\n\n", rule.Name, rule.ServiceName, rule.Path, rule.Status, rule.Method)
		}
	},
}

func Init() {

}

func Execute() {
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(ruleCmd)
	rootCmd.AddCommand(BanCmd)
	rootCmd.AddCommand(UnbanCmd)
	UnbanCmd.Flags().StringVarP(&ip, "ip", "i", "", "ip to unban")
	BanCmd.Flags().StringVarP(&ip, "ip", "i", "", "ip to ban")
	// Rule comand block
	ruleCmd.AddCommand(addCmd)
	ruleCmd.AddCommand(listCmd)
	addCmd.Flags().StringVarP(&name, "name", "n", "", "rule name (required)")
	addCmd.Flags().StringVarP(&service, "service", "s", "", "service name")
	addCmd.Flags().StringVarP(&path, "path", "p", "", "request path")
	addCmd.Flags().StringVarP(&status, "status", "c", "", "HTTP status code")
	addCmd.Flags().StringVarP(&method, "method", "m", "", "HTTP method")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
