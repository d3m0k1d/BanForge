package command

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/judge"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/metrics"
	"github.com/d3m0k1d/BanForge/internal/parser"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"github.com/spf13/cobra"
)

var DaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run BanForge daemon process",
	Run: func(cmd *cobra.Command, args []string) {
		entryCh := make(chan *storage.LogEntry, 1000)
		resultCh := make(chan *storage.LogEntry, 100)
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer stop()
		log := logger.New(false)
		log.Info("Starting BanForge daemon")
		reqDb_w, err := storage.NewRequestsWr()
		if err != nil {
			log.Error("Failed to create request writer", "error", err)
			os.Exit(1)
		}
		reqDb_r, err := storage.NewRequestsRd()
		if err != nil {
			log.Error("Failed to create request reader", "error", err)
			os.Exit(1)
		}
		banDb_r, err := storage.NewBanReader()
		if err != nil {
			log.Error("Failed to create ban reader", "error", err)
			os.Exit(1)
		}
		banDb_w, err := storage.NewBanWriter()
		if err != nil {
			log.Error("Failed to create ban writter", "error", err)
			os.Exit(1)
		}
		defer func() {
			err = banDb_r.Close()
			if err != nil {
				log.Error("Failed to close database connection", "error", err)
			}
			err = banDb_w.Close()
			if err != nil {
				log.Error("Failed to close database connection", "error", err)
			}
		}()
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
		retention, err := config.ParseDurationWithYears(cfg.Storage.RetentionTime)
		if err != nil {
			log.Error("Failed to parse request retention", "error", err)
			os.Exit(1)
		}
		cleanupInterval, err := config.ParseDurationWithYears(cfg.Storage.CleanupInterval)
		if err != nil {
			log.Error("Failed to parse request cleanup interval", "error", err)
			os.Exit(1)
		}
		go func() {
			if err := storage.RequestCleaner(ctx, retention, cleanupInterval, reqDb_w); err != nil {
				log.Error("Request cleanup stopped", "error", err)
			}
		}()
		if cfg.Metrics.Enabled {
			go func() {
				if err := metrics.StartMetricsServer(cfg.Metrics.Port); err != nil {
					log.Error("Failed to start metrics server", "error", err)
				}
			}()
		}
		var b blocker.BlockerEngine
		fw := cfg.Firewall.Name
		b, err = blocker.GetBlocker(fw, cfg.Firewall.Config)
		if err != nil {
			log.Error("Failed to create firewall blocker", "error", err)
			os.Exit(1)
		}
		err = b.Setup(cfg.Firewall.Config)
		if err != nil {
			log.Error("Failed to setup firewall", "error", err)
			os.Exit(1)
		}
		ips, err := banDb_r.RestoreBans()
		if err != nil {
			log.Error("Failed to restore bans", "error", err)
			os.Exit(1)
		}
		for _, ip := range ips {
			err = b.Ban(ip)
			if err != nil {
				log.Error("Failed to ban ip", "ip", ip, "error", err)
			}
		}
		r, err := config.LoadRuleConfig()
		if err != nil {
			log.Error("Failed to load rules", "error", err)
			os.Exit(1)
		}
		j := judge.New(banDb_r, banDb_w, reqDb_r, b, resultCh, entryCh)
		j.LoadRules(r)
		go j.UnbanChecker()

		var parserWg sync.WaitGroup
		var tribunalWg sync.WaitGroup
		var writerWg sync.WaitGroup

		tribunalWg.Add(1)
		go func() {
			defer tribunalWg.Done()
			j.Tribunal()
		}()

		writerWg.Add(1)
		go func() {
			defer writerWg.Done()
			storage.WriteReq(reqDb_w, resultCh)
		}()
		var scanners []*parser.Scanner

		for _, svc := range cfg.Service {
			log.Info(
				"Processing service",
				"name", svc.Name,
				"enabled", svc.Enabled,
				"path", svc.LogPath,
			)

			if !svc.Enabled {
				log.Info("Service disabled, skipping", "name", svc.Name)
				continue
			}

			log.Info("Starting parser for service", "name", svc.Name, "path", svc.LogPath)
			if svc.Logging != "file" && svc.Logging != "journald" {
				log.Error("Invalid logging type", "type", svc.Logging)
				continue
			}

			if svc.Logging == "file" {
				log.Info("Logging to file", "path", svc.LogPath)
				pars, err := parser.NewScannerTail(svc.LogPath)
				if err != nil {
					log.Error("Failed to create scanner", "service", svc.Name, "error", err)
					continue
				}

				scanners = append(scanners, pars)

				go pars.Start()

				parserWg.Add(1)
				go func(p *parser.Scanner, serviceName string) {
					defer parserWg.Done()
					if svc.Name == "nginx" {
						log.Info("Starting nginx parser", "service", serviceName)
						ng := parser.NewNginxParser()
						ng.Parse(p.Events(), entryCh)
					}
					if svc.Name == "ssh" {
						log.Info("Starting ssh parser", "service", serviceName)
						ssh := parser.NewSshdParser()
						ssh.Parse(p.Events(), entryCh)
					}
					if svc.Name == "apache" {
						log.Info("Starting apache parser", "service", serviceName)
						ap := parser.NewApacheParser()
						ap.Parse(p.Events(), entryCh)
					}
				}(pars, svc.Name)
				continue
			}

			if svc.Logging == "journald" {
				log.Info("Logging to journald", "path", svc.LogPath)
				pars, err := parser.NewScannerJournald(svc.LogPath)
				if err != nil {
					log.Error("Failed to create scanner", "service", svc.Name, "error", err)
					continue
				}

				scanners = append(scanners, pars)

				go pars.Start()
				parserWg.Add(1)
				go func(p *parser.Scanner, serviceName string) {
					defer parserWg.Done()
					if svc.Name == "nginx" {
						log.Info("Starting nginx parser", "service", serviceName)
						ng := parser.NewNginxParser()
						ng.Parse(p.Events(), entryCh)

					}
					if svc.Name == "ssh" {
						log.Info("Starting ssh parser", "service", serviceName)
						ssh := parser.NewSshdParser()
						ssh.Parse(p.Events(), entryCh)
					}
					if svc.Name == "apache" {
						log.Info("Starting apache parser", "service", serviceName)
						ap := parser.NewApacheParser()
						ap.Parse(p.Events(), entryCh)
					}

				}(pars, svc.Name)
				continue
			}
		}

		<-ctx.Done()
		log.Info("Shutdown signal received")

		for _, s := range scanners {
			s.Stop()
		}

		parserWg.Wait()
		close(entryCh)
		tribunalWg.Wait()
		close(resultCh)
		writerWg.Wait()

		if err := reqDb_w.Close(); err != nil {
			log.Error("Failed to close request writer database", "error", err)
		}
		if err := reqDb_r.Close(); err != nil {
			log.Error("Failed to close request reader database", "error", err)
		}
		log.Info("BanForge daemon stopped")
	},
}
