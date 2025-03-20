package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	logger            *log.Logger
	DefaultSaveToFile bool
	DefaultLogPath    string
}

/*
 */ /*New creates a new instance of Logger with no flags, saveToFile set false. You can configure default setted values.
 */ /*Example usage:
 */ /*
	func main() {
		log := logger.New()
		log.DefaultLogPath = "logs/app.json"
		startTime := time.Now()

		log.Info(
			"Information message",
			logger.LogOptions{
				StartTime:   startTime,
				User:        "Admin",
				Process:     "MainProcess",
				SaveToFile:  true,
				LogFilePath: "./log.json",
			},
		)
		log.Warn("Warning message", logger.LogOptions{Process: "WorkerProcess", User: "Guest", SaveToFile: true})
		log.Info("Message only")
	}
*/
func New() *Logger {
	return &Logger{
		logger:            log.New(os.Stdout, "", 0),
		DefaultSaveToFile: false,
		DefaultLogPath:    "./log.json",
	}
}

// LogOptions represents optional parameters for logging.
type LogOptions struct {
	StartTime   time.Time // The start time of the process
	Process     string    // The name of the process
	User        string    // The user associated with the log
	SaveToFile  bool      // Whether to save the log to file (default: false)
	LogFilePath string    // Custom path for log file (default: "./log.json")
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
	var opts LogOptions
	if len(options) > 0 {
		opts = options[0]
	}

	// Varsayılan değerleri ayarla
	if opts.LogFilePath == "" {
		opts.LogFilePath = l.DefaultLogPath
	}

	var jsoner logToJSON //var for saving logs to file with json format

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Start()

	startingTime := time.Now()
	s.Prefix = fmt.Sprintf("%s ", startingTime.Format("2006-01-02 15:04:05"))

	logParts := []string{fmt.Sprintf(" [%s]", level)}
	s.Suffix = logParts[0]
	jsoner.Level = level

	time.Sleep(1 * time.Second)
	if opts.Process != "" {
		logParts = append(logParts, opts.Process)
		s.Suffix = strings.Join(logParts, " | ")
		jsoner.Process = opts.Process
	}

	time.Sleep(1 * time.Second)
	if !opts.StartTime.IsZero() {
		currentTime := time.Now()
		duration := fmt.Sprintf("%d ms", currentTime.Sub(opts.StartTime).Milliseconds())

		logParts = append(logParts, duration)
		s.Suffix = strings.Join(logParts, " | ")
		jsoner.Duration = duration
	}

	time.Sleep(1 * time.Second)
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

	if opts.SaveToFile {
		if err := saveToFile(jsoner, opts.LogFilePath); err != nil {
			logParts = append(logParts, "Failed to save log")
			s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " x" + strings.Join(logParts, " | ") + "\n"

			time.Sleep(1 * time.Second)
		} else {
			s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " ✓" + strings.Join(logParts, " | ") + "\n"
		}
	} else {
		s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " ✓" + strings.Join(logParts, " | ") + "\n"
	}

	s.Stop()
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
 */ /* @param "filePath" path to save the log file
 */
func saveToFile(jsoner logToJSON, filePath string) error {
	// Check directory
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create file and initialize array if it doesn't exist
	if !fileExists(filePath) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		_, err = file.WriteString("[\n")
		if err != nil {
			return err
		}
		file.Close()
	}

	// Open file
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// Create JSON data
	jsonData, err := json.MarshalIndent(jsoner, "", "  ")
	if err != nil {
		return err
	}

	if stat.Size() > 2 { // File has at least "[\n"
		// Remove last ']' character
		if err := file.Truncate(stat.Size() - 2); err != nil {
			return err
		}
		// Add comma and newline
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			return err
		}
		if _, err := file.WriteString(",\n"); err != nil {
			return err
		}
	}

	// Append JSON data
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if _, err := file.Write(jsonData); err != nil {
		return err
	}

	// Close array
	_, err = file.WriteString("\n]")
	return err
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
