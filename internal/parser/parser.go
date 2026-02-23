package parser

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/metrics"
)

type Event struct {
	Data string
}

type Scanner struct {
	scanner   *bufio.Scanner
	ch        chan Event
	stopCh    chan struct{}
	logger    *logger.Logger
	cmd       *exec.Cmd
	file      *os.File
	pollDelay time.Duration
}

func validateLogPath(path string) error {
	if path == "" {
		return fmt.Errorf("log path cannot be empty")
	}

	if !filepath.IsAbs(path) {
		return fmt.Errorf("log path must be absolute: %s", path)
	}

	if strings.Contains(path, "..") {
		return fmt.Errorf("log path contains '..': %s", path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("log file does not exist: %s", path)
	}

	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("log path is a symlink: %s", path)
	}

	return nil
}

func validateJournaldUnit(unit string) error {
	if unit == "" {
		return fmt.Errorf("journald unit cannot be empty")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(unit) {
		return fmt.Errorf("invalid journald unit name: %s", unit)
	}

	if strings.HasPrefix(unit, "-") {
		return fmt.Errorf("journald unit cannot start with '-': %s", unit)
	}

	return nil
}

func NewScannerTail(path string) (*Scanner, error) {
	if err := validateLogPath(path); err != nil {
		return nil, fmt.Errorf("invalid log path: %w", err)
	}

	// #nosec G204 - path is validated above via validateLogPath()
	cmd := exec.Command("tail", "-F", "-n", "10", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Scanner{
		scanner:   bufio.NewScanner(stdout),
		ch:        make(chan Event, 100),
		stopCh:    make(chan struct{}),
		logger:    logger.New(false),
		file:      nil,
		cmd:       cmd,
		pollDelay: 100 * time.Millisecond,
	}, nil
}

func NewScannerJournald(unit string) (*Scanner, error) {
	if err := validateJournaldUnit(unit); err != nil {
		return nil, fmt.Errorf("invalid journald unit: %w", err)
	}

	// #nosec G204 - unit is validated above via validateJournaldUnit()
	cmd := exec.Command("journalctl", "-u", unit, "-f", "-n", "0", "-o", "short", "--no-pager")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Scanner{
		scanner:   bufio.NewScanner(stdout),
		ch:        make(chan Event, 100),
		stopCh:    make(chan struct{}),
		logger:    logger.New(false),
		cmd:       cmd,
		file:      nil,
		pollDelay: 100 * time.Millisecond,
	}, nil
}

func (s *Scanner) Start() {
	s.logger.Info("Scanner started")

	go func() {
		for {
			select {
			case <-s.stopCh:
				s.logger.Info("Scanner stopped")
				return

			default:
				if s.scanner.Scan() {
					metrics.IncScannerEvent("scanner")
					s.ch <- Event{
						Data: s.scanner.Text(),
					}
					s.logger.Info("Scanner event", "data", s.scanner.Text())
				} else {
					if err := s.scanner.Err(); err != nil {
						s.logger.Error("Scanner error")
						metrics.IncError()
						return
					}
				}
			}
		}
	}()
}

func (s *Scanner) Stop() {
	close(s.stopCh)

	if s.cmd != nil && s.cmd.Process != nil {
		s.logger.Info("Stopping process", "pid", s.cmd.Process.Pid)
		err := s.cmd.Process.Kill()
		if err != nil {
			s.logger.Error("Failed to kill process", "err", err)
		}
		err = s.cmd.Wait()
		if err != nil {
			s.logger.Error("Failed to wait process", "err", err)
		}

	}

	if s.file != nil {
		if err := s.file.Close(); err != nil {
			s.logger.Error("Failed to close file", "err", err)
		}
	}
	time.Sleep(150 * time.Millisecond)
	close(s.ch)
}

func (s *Scanner) Events() <-chan Event {
	return s.ch
}
