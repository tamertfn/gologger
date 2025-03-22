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

/*
 */ /*FileLog is a boolean or string that determines whether to save logs to a file.
 */ /*false for no file logging, true for default file logging, filepath (string) for custom file logging.
 */
type Logger struct {
	logger  *log.Logger
	FileLog any // default file logging setting
}

/*
 */ /*New creates a new instance of Logger with no flags, FileLog set false.
 */ /*You can configure use .SetDefaultFileLog() method to change default file logging setting.
 */ /*Example usage:
 */ /*
func main() {
	// Create a new logger
	log := logger.New()

	// 1. Test with default settings (no file logging)
	log.Info("Test 1: Default settings - no file logging")

	// 2. Test with default set to true
	log.SetDefaultFileLog(true)
	log.Info("Test 2: Default set to true - will log to default path")

	// 3. Test with custom default path
	log.SetDefaultFileLog("./logs/custom_default.json")
	log.Info("Test 3: Default set to custom path")

	// 4. Test with single override
	log.Info("Test 4: Override default path for single log", logger.LogOptions{
		FileLog: "./logs/single_override.json",
	})

	// 5. Test with all features
	log.Info("Test 5: Full features test", logger.LogOptions{
		StartTime: time.Now(),
		Process:   "MainProcess",
		User:      "TestUser",
		FileLog:   "./logs/full_test.json",
	})

	// 6. Test error handling (INFO -> WARN conversion)
	log.Info("Test 6: Error handling test", logger.LogOptions{
		FileLog: "/invalid/path/test.json", // This will cause an error
	})

	// 7. Test with default set to false and override
	log.SetDefaultFileLog(false)
	log.Info("Test 7: Default false with override", logger.LogOptions{
		FileLog: true, // This will log to default.json
	})

	// 8. Process tracking example
	startTime := time.Now()
	time.Sleep(200 * time.Millisecond) // Simulated process
	log.Info("Test 8: Process duration test", logger.LogOptions{
		StartTime: startTime,
		Process:   "SlowProcess",
		FileLog:   "./logs/process_test.json",
	})
}
*/
func New() *Logger {
	return &Logger{
		logger:  log.New(os.Stdout, "", 0),
		FileLog: false,
	}
}

// LogOptions represents optional parameters for logging.
type LogOptions struct {
	StartTime time.Time
	Process   string
	User      string
	FileLog   any // can be bool or string
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

	// Use provided FileLog from options if set, otherwise use default
	fileLogSetting := l.FileLog
	if opts.FileLog != nil {
		fileLogSetting = opts.FileLog
	}

	var jsoner logToJSON //var for saving logs to file with json format

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Start()

	startingTime := time.Now()
	s.Prefix = fmt.Sprintf("%s ", startingTime.Format("2006-01-02 15:04:05"))

	logParts := []string{fmt.Sprintf(" [%s]", level)}
	s.Suffix = logParts[0]
	jsoner.Level = level

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

	// Handle file logging based on FileLog type
	shouldSaveToFile, filePath, err := checkSaveLogOption(fileLogSetting)
	if err != nil {
		errorMsg := fmt.Sprintf("Error: %v", err)
		logParts = append(logParts, errorMsg)
		s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " x" + strings.Join(logParts, " | ") + "\n"
	}

	if shouldSaveToFile {
		if err := saveToFile(jsoner, filePath); err != nil {
			if level == INFO {
				logParts[0] = fmt.Sprintf(" [%s]", WARN)
				jsoner.Level = WARN
			}
			errorMsg := fmt.Sprintf("Error: %v", err)
			logParts = append(logParts, errorMsg)
			s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " x" + strings.Join(logParts, " | ") + "\n"
		} else {
			logParts = append(logParts, "Log saved!")
			s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " ✓" + strings.Join(logParts, " | ") + "\n"
		}
	} else {
		s.FinalMSG = closingTime.Format("2006-01-02 15:04:05") + " ✓" + strings.Join(logParts, " | ") + "\n"
	}

	s.Stop()
}

/*
 */ /*INFO Level Logging Handler
 */ /*Switches to WARN if save to log fails.
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
		return fmt.Errorf("directory creation failed, %v", err)
	}

	// Create file and initialize array if it doesn't exist
	if !fileExists(filePath) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("file creation failed, %v", err)
		}
		defer file.Close()
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("file open failed, %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file stat failed, %v", err)
	}

	jsonData, err := json.MarshalIndent(jsoner, "  ", "  ")
	if err != nil {
		return fmt.Errorf("JSON marshaling failed, %v", err)
	}

	// Eğer dosya boşsa
	if stat.Size() == 0 {
		// İlk kayıt için array başlat
		if _, err := file.WriteString("[\n  " + string(jsonData) + "\n]"); err != nil {
			return fmt.Errorf("initial write failed, %v", err)
		}
		return nil
	}

	// Dosya sonundaki ']' karakterini sil
	if err := file.Truncate(stat.Size() - 1); err != nil {
		return fmt.Errorf("truncate failed, %v", err)
	}

	// Dosya sonuna git
	if _, err := file.Seek(-1, io.SeekEnd); err != nil {
		return fmt.Errorf("seek failed, %v", err)
	}

	// Yeni kaydı ekle
	if _, err := file.WriteString(",\n  " + string(jsonData) + "\n]"); err != nil {
		return fmt.Errorf("append failed, %v", err)
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// SetDefaultFileLog allows changing the default file logging behavior
func (l *Logger) SetDefaultFileLog(fileLog any) {
	l.FileLog = fileLog
}

// checkLogOption checks FileLog setting and returns if should save and where to save
func checkSaveLogOption(fileLogSetting any) (shouldSave bool, path string, err error) {
	path = "./log.json" // default path

	switch v := fileLogSetting.(type) {
	case bool:
		shouldSave = v
	case string:
		if v != "" {
			shouldSave = true
			path = v
		}
	default:
		return false, "", fmt.Errorf("log will not be saved, invalid file log setting type: %T, expected bool or string", fileLogSetting)
	}

	return shouldSave, path, nil
}

// Test runs all test cases for the logger
func (l *Logger) Test() {
	// 1. Test with default settings (no file logging)
	l.Info("Test 1: Default settings - no file logging")

	// 2. Test with custom default path
	l.SetDefaultFileLog("./logs/custom_default.json")
	l.Info("Test 3: Default set to custom path")

	// 3. Test with single override
	l.Info("Test 4: Override default path for single log", LogOptions{
		FileLog: "./logs/single_override.json",
	})

	// 4. Test with all features
	l.Info("Test 5: Full features test", LogOptions{
		StartTime: time.Now(),
		Process:   "MainProcess",
		User:      "TestUser",
		FileLog:   "./logs/full_test.json",
	})

	// 5. Test error handling (INFO -> WARN conversion)
	l.Info("Test 6: Error handling test", LogOptions{
		FileLog: "/invalid/path/test.json", // This will cause an error
	})

	// 6. Test with default set to false and override
	l.SetDefaultFileLog(false)
	l.Info("Test 7: Default false with override", LogOptions{
		FileLog: true, // This will log to default.json
	})

	// 7. Process tracking example
	startTime := time.Now()
	time.Sleep(200 * time.Millisecond) // Simulated process
	l.Info("Test 8: Process duration test", LogOptions{
		StartTime: startTime,
		Process:   "SlowProcess",
		FileLog:   "./logs/process_test.json",
	})
}

func (l *Logger) TestWithPanic() {
	// 1. Test with default settings (no file logging)
	l.Info("Test 1: Default settings - no file logging")

	// 2. Test with custom default path
	l.SetDefaultFileLog("./logs/custom_default.json")
	l.Info("Test 3: Default set to custom path")

	// 3. Test with single override
	l.Info("Test 4: Override default path for single log", LogOptions{
		FileLog: "./logs/single_override.json",
	})

	// 4. Test with all features
	l.Info("Test 5: Full features test", LogOptions{
		StartTime: time.Now(),
		Process:   "MainProcess",
		User:      "TestUser",
		FileLog:   "./logs/full_test.json",
	})

	// 5. Test error handling (INFO -> WARN conversion)
	l.Info("Test 6: Error handling test", LogOptions{
		FileLog: "/invalid/path/test.json", // This will cause an error
	})

	// 6. Test with default set to false and override
	l.SetDefaultFileLog(false)
	l.Info("Test 7: Default false with override", LogOptions{
		FileLog: true, // This will log to default.json
	})

	// 7. Process tracking example
	startTime := time.Now()
	time.Sleep(200 * time.Millisecond) // Simulated process
	l.Info("Test 8: Process duration test", LogOptions{
		StartTime: startTime,
		Process:   "SlowProcess",
		FileLog:   "./logs/process_test.json",
	})

	// 8. Test Panic with recovery
	defer func() {
		if r := recover(); r != nil {
			l.Info("Recovered from panic", LogOptions{
				Process: "PanicRecovery",
				FileLog: "./logs/recovery.json",
			})
		}
	}()

	// Test Panic with options
	l.Panic("Test 9: Panic test with options", LogOptions{
		Process: "PanicProcess",
		User:    "TestUser",
		FileLog: "./logs/panic_test.json",
	})
}

func (l *Logger) TestWithFatal() {
	// 1. Test with default settings (no file logging)
	l.Info("Test 1: Default settings - no file logging")

	// 2. Test with custom default path
	l.SetDefaultFileLog("./logs/custom_default.json")
	l.Info("Test 3: Default set to custom path")

	// 3. Test with single override
	l.Info("Test 4: Override default path for single log", LogOptions{
		FileLog: "./logs/single_override.json",
	})

	// 4. Test with all features
	l.Info("Test 5: Full features test", LogOptions{
		StartTime: time.Now(),
		Process:   "MainProcess",
		User:      "TestUser",
		FileLog:   "./logs/full_test.json",
	})

	// 5. Test error handling (INFO -> WARN conversion)
	l.Info("Test 6: Error handling test", LogOptions{
		FileLog: "/invalid/path/test.json", // This will cause an error
	})

	// 6. Test with default set to false and override
	l.SetDefaultFileLog(false)
	l.Info("Test 7: Default false with override", LogOptions{
		FileLog: true, // This will log to default.json
	})

	// 7. Process tracking example
	startTime := time.Now()
	time.Sleep(200 * time.Millisecond) // Simulated process
	l.Info("Test 8: Process duration test", LogOptions{
		StartTime: startTime,
		Process:   "SlowProcess",
		FileLog:   "./logs/process_test.json",
	})

	l.Fatal("Test 10: Fatal test with options", LogOptions{
		Process: "FatalProcess",
		User:    "TestUser",
		FileLog: "./logs/fatal_test.json",
	})
}
