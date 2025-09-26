package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

// NewLogger creates a new logger with output to both file and stdout/stderr
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Info → stdout + file
	infoWriter := io.MultiWriter(os.Stdout, file)
	// Error → stderr + file
	errorWriter := io.MultiWriter(os.Stderr, file)

	return &Logger{
		infoLogger:  log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		file:        file,
	}, nil
}

// Info logs informational messages
func (l *Logger) Info(msg string) {
	l.infoLogger.Printf("%s", msg)
}

// Error logs error messages (with optional error value)
func (l *Logger) Error(msg string, err error) {
	if err != nil {
		l.errorLogger.Printf("%s: %v", msg, err)
	} else {
		l.errorLogger.Printf("%s", msg)
	}
}

// Close closes the log file
func (l *Logger) Close() error {
	return l.file.Close()
}
