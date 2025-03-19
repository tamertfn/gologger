package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	INFO  = "INFO"
	WARN  = "WARN"
	FATAL = "FATAL"
	PANIC = "PANIC"
)

// Logger represents a logging system.
type Logger struct {
	logger *log.Logger
}

// New creates a new instance of Logger with Ldate|Ltime flags.
func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// LogOptions represents optional parameters for logging.
type LogOptions struct {
	StartTime time.Time // The start time of the process
	Process   string    // The name of the process
	User      string    // The user associated with the log
}

/*
logMessage logs a message with a specific level.
Parameters:
  - level: The log level (INFO, WARN, FATAL, PANIC)
  - message: The log message
  - options: Optional parameters such as StartTime, Process, and User
*/
func (l *Logger) logMessage(level, message string, options ...LogOptions) {
	logParts := []string{fmt.Sprintf("[%s]", level)}

	var opts LogOptions
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Process != "" {
		logParts = append(logParts, opts.Process)
	}

	if !opts.StartTime.IsZero() {
		currentTime := time.Now()
		duration := fmt.Sprintf("%d ms", currentTime.Sub(opts.StartTime).Milliseconds())
		logParts = append(logParts, duration)
	}

	if opts.User != "" {
		logParts = append(logParts, opts.User)
	}

	logParts = append(logParts, message)
	l.logger.Println("|", strings.Join(logParts, " | "))
}

/*
Info logs an informational message.
Parameters:
  - message: The log message
  - options: Optional parameters such as StartTime, Process, and User
*/
func (l *Logger) Info(message string, options ...LogOptions) {
	l.logMessage(INFO, message, options...)
}

/*
Warn logs a warning message.
Parameters:
  - message: The log message
  - options: Optional parameters such as StartTime, Process, and User
*/
func (l *Logger) Warn(message string, options ...LogOptions) {
	l.logMessage(WARN, message, options...)
}

/*
Fatal logs a fatal error message and exits the program with status code of 1.
Parameters:
  - message: The log message
  - options: Optional parameters such as StartTime, Process, and User
*/
func (l *Logger) Fatal(message string, options ...LogOptions) {
	l.logMessage(FATAL, message, options...)
	os.Exit(1)
}

/*
Panic logs a fatal error message and exits the program with status code of 2.
Parameters:
  - message: The log message
  - options: Optional parameters such as StartTime, Process, and User
*/
func (l *Logger) Panic(message string, options ...LogOptions) {
	l.logMessage(PANIC, message, options...)
	panic(message)
}
