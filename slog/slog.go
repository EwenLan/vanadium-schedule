package slog

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

type slog struct {
	log.Logger
	disableStandardLogOutputOption bool
	runningMode                    bool
}

// Init
func (l *slog) Init(filepath string) {
	directory := path.Dir(filepath)
	err0 := os.MkdirAll(directory, logPermit)
	if err0 != nil {
		fmt.Printf("make directory failed, err = %s", err0)
		return
	}
	logFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logPermit)
	if err != nil {
		fmt.Printf("open log file failed, err: %s\n", err)
		return
	}
	l.SetOutput(logFile)
	l.runningMode = true
}

func (l *slog) getSingleLogRow(level, format string, args ...interface{}) string {
	_, file, line, ok := runtime.Caller(outputDepth)
	if !ok || file == "" {
		file = "???"
		line = 999
	}
	prefix := fmt.Sprintf("[%s][%s][%s:%d]", level, time.Now().Format(time.RFC3339), path.Base(file), line)
	return fmt.Sprintf("["+prefix+"["+format+"]]", args...)
}

func (l *slog) output(level string, format string, args ...interface{}) {
	row := l.getSingleLogRow(level, format, args...)
	if l.runningMode {
		l.Output(0, row)
	}
	if !l.runningMode || !l.disableStandardLogOutputOption {
		fmt.Println(row)
	}
}

func (l *slog) SetDisableStandardLogOutput(option bool) {
	l.disableStandardLogOutputOption = option
	l.Debugf("set standard log option = %t", option)
}

// Debugf
func (l *slog) Debugf(format string, args ...interface{}) {
	l.output(debug, format, args...)
}

// Infof
func (l *slog) Infof(format string, args ...interface{}) {
	l.output(info, format, args...)
}

// Errorf
func (l *slog) Errorf(format string, args ...interface{}) {
	l.output(error, format, args...)
}
