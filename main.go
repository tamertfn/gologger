package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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

/*
New creates a new instance of Logger with no flags.

Example usage:

	func main() {
		log := logger.New()
		startTime := time.Now()

		log.Info(
			"Information message",
			logger.LogOptions{
				StartTime: startTime,
				User:      "Admin",
				Process:   "MainProcess",
			},
		)
		log.Warn("Warning message", logger.LogOptions{User: "Guest", Process: "WorkerProcess"})
		log.Info("Message only")
	}
*/
func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
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

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Start()

	startingTime := time.Now()
	s.Prefix = fmt.Sprintf("%s ", startingTime.Format("2006-01-02 15:04:05"))

	logParts := []string{fmt.Sprintf(" [%s]", level)}
	s.Suffix = logParts[0]

	var opts LogOptions
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Process != "" {
		logParts = append(logParts, opts.Process)
		s.Suffix = strings.Join(logParts, " | ")
	}

	if !opts.StartTime.IsZero() {
		currentTime := time.Now()
		duration := fmt.Sprintf("%d ms", currentTime.Sub(opts.StartTime).Milliseconds())
		logParts = append(logParts, duration)
		s.Suffix = strings.Join(logParts, " | ")
	}

	if opts.User != "" {
		logParts = append(logParts, opts.User)
		s.Suffix = strings.Join(logParts, " | ")
	}

	logParts = append(logParts, message)
	s.Suffix = strings.Join(logParts, " | ")

	closingTime := time.Now()
	s.FinalMSG = fmt.Sprintf("%s ", closingTime.Format("2006-01-02 15:04:05")) + "✓" + strings.Join(logParts, " | ") + "\n"

	s.Stop()
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
