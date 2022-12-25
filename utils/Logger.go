package utils

import (
	"io"
	"log"
)

// Logger manage concurrent logs
type Logger struct {
	LogChan chan interface{}
	Running bool
	logger  *log.Logger
}

// NewLogger creates a *Logger
func NewLogger(writer io.Writer) *Logger {
	return &Logger{
		LogChan: make(chan interface{}),
		logger:  log.New(writer, "", log.Flags()),
	}
}

// Logs go function for logging
func (l *Logger) Logs() {
	l.Running = true
	for l.Running {
		select {
		case message := <-l.LogChan:
			log.Printf("%v", message)
		}
	}
}
