package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type Level int

//go:generate stringer -type=Level
const (
	DEBUG Level = iota + 1
	INFO
	ERROR
)

//nolint:lll
//go:generate mockgen -destination=../mocks/mock_logger.go -package=mocks github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger Logger
type Logger interface {
	ResetWriters()
	AddWriter(writer SyncWriter)
	AddLogFile(logFile string) (*os.File, error)
	Save(lvl Level, msg string, params ...Parameter)
	Debug(msg string, params ...Parameter)
	Info(msg string, params ...Parameter)
	Error(msg string, params ...Parameter)
}

type SyncWriter interface {
	io.Writer
	Sync() error
}

type logg struct {
	minLevel Level
	writers  []SyncWriter
}

func New(options ...Option) Logger {
	logg := &logg{
		minLevel: INFO,
	}
	logg.AddWriter(os.Stderr)
	logg.fillOptions(options)
	return logg
}

func (l *logg) ResetWriters() {
	l.writers = nil
}

func (l *logg) AddWriter(writer SyncWriter) {
	l.writers = append(l.writers, writer)
}

func (l *logg) AddLogFile(logFile string) (*os.File, error) {
	dir := path.Dir(logFile)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, err
		}
	}
	if !os.IsNotExist(err) && err != nil {
		return nil, err
	}
	fp, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	l.AddWriter(fp)
	return fp, nil
}

func (l *logg) Save(lvl Level, msg string, params ...Parameter) {
	if lvl < l.minLevel {
		return
	}
	for _, fp := range l.writers {
		fmt.Fprintf(
			fp,
			"%s\t%s\t%s\t[%s]\n",
			time.Now().Format(time.RFC3339),
			lvl.String(),
			msg,
			l.joinParams(params),
		)
		fp.Sync()
	}
}

func (l *logg) Debug(msg string, params ...Parameter) {
	l.Save(DEBUG, msg, params...)
}

func (l *logg) Info(msg string, params ...Parameter) {
	l.Save(INFO, msg, params...)
}

func (l *logg) Error(msg string, params ...Parameter) {
	l.Save(ERROR, msg, params...)
}

func (l *logg) fillOptions(options []Option) {
	for _, option := range options {
		opt, ok := option.(*OptionMinLevel)
		if ok {
			l.minLevel = opt.getOption().(Level)
		}
	}
}

func (l *logg) joinParams(params []Parameter) string {
	strs := make([]string, len(params))
	for i, param := range params {
		strs[i] = param.GetKeyValue()
	}
	return strings.Join(strs, ", ")
}
