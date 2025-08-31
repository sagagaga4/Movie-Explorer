package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger //For info massages
	errorLogger *log.Logger //For error massages
	file        *os.File
}

// NewLogger creates a new logger with output to both file and stdout
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		file:        file,
	}, nil
}

// Info logs informational messages to stdout
func (l *Logger) Info(msg string) {
	l.infoLogger.Printf("%s", msg)
}

// Error logs error messages to file
func (l *Logger) Error(msg string, err error) {
	l.errorLogger.Printf("%s: %v", msg, err)
}

// Close closes the log file
func (l *Logger) Close() {
	l.file.Close()
}
