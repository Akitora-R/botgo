package log

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

var _ Logger = (*consoleLogger)(nil)

type consoleLogger struct {
	Level Level
}

var levelS = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}

type Level int

func (l Level) String() string {
	return levelS[l]
}

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (l *consoleLogger) evalLevel(level Level) bool {
	return l.Level <= level
}

func (l *consoleLogger) Debug(v ...interface{}) {
	l.output(DebugLevel, fmt.Sprint(v...))
}

func (l *consoleLogger) Info(v ...interface{}) {
	l.output(InfoLevel, fmt.Sprint(v...))
}

func (l *consoleLogger) Warn(v ...interface{}) {
	l.output(WarnLevel, fmt.Sprint(v...))
}

func (l *consoleLogger) Error(v ...interface{}) {
	l.output(ErrorLevel, fmt.Sprint(v...))
}

func (l *consoleLogger) Debugf(format string, v ...interface{}) {
	l.output(DebugLevel, fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Infof(format string, v ...interface{}) {
	l.output(InfoLevel, fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Warnf(format string, v ...interface{}) {
	l.output(WarnLevel, fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Errorf(format string, v ...interface{}) {
	l.output(ErrorLevel, fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Sync() error {
	return nil
}

func (l *consoleLogger) output(level Level, v ...interface{}) {
	if !l.evalLevel(level) {
		return
	}
	pc, file, line, _ := runtime.Caller(3)
	file = filepath.Base(file)
	funcName := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")

	//logFormat := "%-7v %s %-17s: %s " + fmt.Sprint(v...) + "\n"
	//date := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Printf(logFormat, fmt.Sprintf("[%v]", level), date, fmt.Sprintf("[%s:%d]", file, line), funcName)
	switch level {
	case DebugLevel:
		slog.Debug(fmt.Sprint(v...), "pos", fmt.Sprintf("[%s:%d]", file, line), "func", funcName)
	case InfoLevel:
		slog.Info(fmt.Sprint(v...), "pos", fmt.Sprintf("[%s:%d]", file, line), "func", funcName)
	case WarnLevel:
		slog.Warn(fmt.Sprint(v...), "pos", fmt.Sprintf("[%s:%d]", file, line), "func", funcName)
	case ErrorLevel:
		slog.Error(fmt.Sprint(v...), "pos", fmt.Sprintf("[%s:%d]", file, line), "func", funcName)
	}
}
