package command

import (
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

var DaemonCmd = &cobra.Command{
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
