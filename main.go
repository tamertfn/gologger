package logger

import (
	"encoding/json"
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
 */ /*New creates a new instance of Logger with no flags.
 */ /*Example usage:
 */ /*
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

type logToJSON struct {
	ClosingTime string `json:"time"`
	Level       string `json:"level"`
	Process     string `json:"process,omitempty"`
	Duration    string `json:"duration,omitempty"`
	User        string `json:"user,omitempty"`
	Message     string `json:"message"`
}

/*
 */ /*logMassage Handler
 */ /* @param "level" the log level (INFO, WARN, FATAL, PANIC)
 */ /* @param "message" message to be logged
 */ /* @param "options" optional parameters such as StartTime, Process, and User
 */
func (l *Logger) logMessage(level, message string, options ...LogOptions) {

	var jsoner logToJSON //var for saving logs to file with json format

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Start()

	startingTime := time.Now()
	s.Prefix = fmt.Sprintf("%s ", startingTime.Format("2006-01-02 15:04:05"))

	logParts := []string{fmt.Sprintf(" [%s]", level)}
	s.Suffix = logParts[0]
	jsoner.Level = level

	var opts LogOptions
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Process != "" {
		logParts = append(logParts, opts.Process)
		s.Suffix = strings.Join(logParts, " | ")
		jsoner.Process = opts.Process
	}

	if !opts.StartTime.IsZero() {
		currentTime := time.Now()
		duration := fmt.Sprintf("%d ms", currentTime.Sub(opts.StartTime).Milliseconds())

		logParts = append(logParts, duration)
		s.Suffix = strings.Join(logParts, " | ")
		jsoner.Duration = duration
	}

	if opts.User != "" {
		logParts = append(logParts, opts.User)
		s.Suffix = strings.Join(logParts, " | ")
		jsoner.User = opts.User
	}

	logParts = append(logParts, message)
	s.Suffix = strings.Join(logParts, " | ")
	jsoner.Message = message

	closingTime := time.Now()
	jsoner.ClosingTime = closingTime.Format("2006-01-02 15:04:05")
	s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " ✓" + strings.Join(logParts, " | ") + "\n"

	s.Stop()

	saveToFile(jsoner)

}

/*
 */ /*INFO Level Logging Handler
 */ /* @param "message" message to be logged
 */ /* @param "options" optional parameters such as StartTime, Process, and User
 */
func (l *Logger) Info(message string, options ...LogOptions) {
	l.logMessage(INFO, message, options...)
}

/*
 */ /*WARN Level Logging Handler
 */ /* @param "message" message to be logged
 */ /* @param "options" optional parameters such as StartTime, Process, and User
 */
func (l *Logger) Warn(message string, options ...LogOptions) {
	l.logMessage(WARN, message, options...)
}

/*
 */ /*FATAL Level Logging Handler
 */ /* @param "message" message to be logged
 */ /* @param "options" optional parameters such as StartTime, Process, and User
 */ /* exits with os.exit(1)
 */
func (l *Logger) Fatal(message string, options ...LogOptions) {
	l.logMessage(FATAL, message, options...)
	os.Exit(1)
}

/*
 */ /*PANIC Level Logging Handler
 */ /* @param "message" message to be logged
 */ /* @param "options" optional parameters such as StartTime, Process, and User
 */ /* exits with panic
 */
func (l *Logger) Panic(message string, options ...LogOptions) {
	l.logMessage(PANIC, message, options...)
	panic(message)
}

/*
 */ /*Saves logs to file in JSON format
 */ /* @param "jsoner" logToJSON struct
 */
func saveToFile(jsoner logToJSON) error {
	file, err := os.OpenFile("log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(jsoner)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(jsonData) + "\n")
	return err
}
