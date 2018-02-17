package core

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/fatih/color"
)

// customize logging
// i know what your saying right now, YAGNI, but i can't help myself (‾⌣‾)♉
// i could probably write all these in a better way by implementing a custom io.Writer
// for each case (success, info, error) and wrap the text in a color of choice and specify
// a prefix on an individual basis, but i don't feel like it. But i probably could have
// done that during the time it took me to write this. #designdecisions
type Logger struct {
	success *log.Logger
	info    *log.Logger
	error   *log.Logger
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{
		success: log.New(w, "", 0), // "you *really* don't know what this zero means?" -- ex-coworker
		info:    log.New(w, "", 0),
		error:   log.New(w, "", log.Lshortfile),
	}

}

func (l *Logger) Success(format string, v ...interface{}) {
	t := time.Now().Format(time.RFC3339)
	s := fmt.Sprintf(t+": "+format, v...)

	l.success.Print(color.GreenString(s))
}

func (l *Logger) Info(format string, v ...interface{}) {
	t := time.Now().Format(time.RFC3339)
	s := fmt.Sprintf(t+": "+format, v...)

	l.info.Print(color.BlueString(s))
}

func (l *Logger) Error(format string, v ...interface{}) {
	t := time.Now().Format(time.RFC3339)
	s := fmt.Sprintf(format, v...)
	l.error.SetPrefix(color.RedString(t + ": "))

	l.error.Print(color.RedString(s))
}
