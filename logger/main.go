package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.warningLogger.Printf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

func SetupLogger() *Logger {
	return &Logger{
		infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
