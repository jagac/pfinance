package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogWriter is a custom log writer that writes logs to both stdout and a log file.
// It rotates the log file daily based on the current date.
type LogWriter struct {
	stdout     *os.File // Standard output (stdout) file
	file       *os.File // Log file for writing logs
	logDir     string   // Directory for storing log files
	currentDay string   // Tracks the current day for log rotation
}

// Write writes the log data to both stdout and the log file.
// It checks for log file rotation if the date changes.
func (t *LogWriter) Write(p []byte) (n int, err error) {

	currentTime := time.Now().Format("2006-01-02")

	if currentTime != t.currentDay {
		if t.file != nil {
			err := t.file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error closing log file: %v\n", err)
			}
		}

		t.currentDay = currentTime
		logFileName := fmt.Sprintf("%s/logs_%s.log", t.logDir, currentTime)

		if !strings.HasPrefix(logFileName, t.logDir+"/") {
			return 0, fmt.Errorf("invalid log file path")
		}

		t.file, err = os.OpenFile(filepath.Clean(logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return 0, err
		}
	}

	n, err = t.stdout.Write(p)
	if err != nil {
		return n, err
	}

	n, err = t.file.Write(p)
	return n, err
}

// NewLogger creates a new logger that writes logs to both stdout and a daily log file.
// The log file is stored in the directory specified by the LOG_DIR environment variable.
func NewLogger() (*slog.Logger, error) {
	logDir := "logs"

	if err := os.MkdirAll(logDir, 0750); err != nil {
		return nil, err
	}

	currentTime := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s/logs_%s.log", logDir, currentTime)

	if !strings.HasPrefix(logFileName, logDir+"/") {
		return nil, fmt.Errorf("invalid log file path")
	}

	file, err := os.OpenFile(filepath.Clean(logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	writer := &LogWriter{
		stdout:     os.Stdout,
		file:       file,
		logDir:     logDir,
		currentDay: currentTime,
	}

	h := slog.NewJSONHandler(writer, nil)
	logger := slog.New(h)

	return logger, nil
}
