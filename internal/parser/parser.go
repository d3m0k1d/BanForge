package parser

import (
	"bufio"
	"os"
	"time"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type Event struct {
	Data string
}

type Scanner struct {
	scanner   *bufio.Scanner
	ch        chan Event
	stopCh    chan struct{}
	logger    *logger.Logger
	file      *os.File
	pollDelay time.Duration
}

func NewScanner(path string) (*Scanner, error) {
	file, err := os.Open(
		path,
	) // #nosec G304 -- admin tool, runs as root, path controlled by operator
	if err != nil {
		return nil, err
	}

	return &Scanner{
		scanner:   bufio.NewScanner(file),
		ch:        make(chan Event, 100),
		stopCh:    make(chan struct{}),
		logger:    logger.New(false),
		file:      file,
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
					s.ch <- Event{
						Data: s.scanner.Text(),
					}
					s.logger.Info("Scanner event", "data", s.scanner.Text())
				} else {
					if err := s.scanner.Err(); err != nil {
						s.logger.Error("Scanner error")
						return
					}
					time.Sleep(s.pollDelay)
				}
			}
		}
	}()
}

func (s *Scanner) Stop() {
	close(s.stopCh)
	time.Sleep(150 * time.Millisecond)
	err := s.file.Close()
	if err != nil {
		s.logger.Error("Failed to close file")
	}
	close(s.ch)
}

func (s *Scanner) Events() <-chan Event {
	return s.ch
}
