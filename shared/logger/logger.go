package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
)

// Logger defines methods for logging messages and errors with a trace ID, method context, and log type.
type Logger interface {
	// Log records a message with the given trace ID, method name, and log type.
	Log(traceID, method, message string, logType LogType)
	// LogAndReturnError logs the provided error with trace ID, method name, and log type,
	// then returns the same error for further handling.
	LogAndReturnError(traceID, method string, err error, logType LogType) error
}

// NewLogger creates a new Logger instance for the specified class.
//
// It requires the file path (file) where logs will be written, and the class name (class)
// indicating the context where this logger will be used.
func NewLogger(file, class string) Logger {
	if config == nil {
		LoadDefaultLogConfig()
	}
	return &logger{file: file, class: class}
}

type logger struct {
	file  string
	class string
}

func (l *logger) Log(traceID, method, message string, logType LogType) {
	if config[logType].ToConsole {
		l.logToConsole(l.formString(traceID, method, message, logType))
	}

	if config[logType].ToFile {
		l.logToFile(l.formString(traceID, method, message, logType))
	}
}
func (l *logger) LogAndReturnError(traceID, method string, err error, logType LogType) error {
	l.Log(traceID, "Add", err.Error(), logType)
	return err
}

func (l *logger) formString(traceID, method, message string, logType LogType) string {
	return fmt.Sprintf("Type: %s, traceID: %s, time: %s, Struct: %s, Method: %s ======> %s",
		logType, traceID, time.Now().Format(db.TimeFormat), l.class, method, message,
	)
}

func (l *logger) logToFile(message string) {
	f, err := os.OpenFile(l.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.logToConsole(l.formString(
			"", "logToFile", fmt.Sprintf("failed to open log file %s: %v\n", l.file, err), LoggerError,
		))
		return
	}
	defer f.Close()

	if _, err := f.WriteString(message + "\n"); err != nil {
		fmt.Printf("failed to write log to file %s: %v\n", l.file, err)
	}
}

func (l *logger) logToConsole(message string) {
	fmt.Println(message)
}
