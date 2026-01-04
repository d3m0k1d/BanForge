package parser

import (
	"bufio"
	"github.com/d3m0k1d/BanForge/local/logger"
	"os"
)

type Event struct {
	Data string
}

type Scaner struct {
	scanner *bufio.Scanner
	ch      chan Event
}

func CreateScaner(path string) *Scaner {
	log := logger.New(false)
	file, err := os.Open(path)
	if err != nil {
		log.Error(err.Error())
	}
	return &Scaner{
		scanner: bufio.NewScanner(file),
		ch:      make(chan Event),
	}
}
