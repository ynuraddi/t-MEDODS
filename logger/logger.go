package logger

import (
	"fmt"
	"io"
	"log"

	"github.com/ynuraddi/t-medods/config"
)

type Logger interface {
	Error(msg string, err error)
	Info(msg string)
}

const (
	ERR = iota
	INF
)

type logger struct {
	err *log.Logger
	inf *log.Logger

	level int
}

func NewLogger(config *config.Config, out io.Writer) *logger {
	logger := logger{
		level: config.LogLevel,
	}

	lerr := log.New(out, "[ERR]\t", log.Ltime|log.Lshortfile)
	linf := log.New(out, "[INF]\t", log.Ltime)

	logger.err = lerr
	logger.inf = linf

	return &logger
}

func (l *logger) Error(msg string, err error) {
	if l.level < ERR {
		return
	}

	l.err.Output(2, fmt.Sprintln(msg, err))
}

func (l *logger) Info(msg string) {
	if l.level < INF {
		return
	}

	l.inf.Println(msg)
}
