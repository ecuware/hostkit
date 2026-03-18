package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles logging operations
type Manager struct {
	logDir   string
	maxSize  int64 // Maximum log file size in bytes
	maxFiles int   // Maximum number of log files to keep
}

// NewManager creates a new log manager
func NewManager(logDir string) *Manager {
	if logDir == "" {
		logDir = "/var/log/hostkit"
	}

	return &Manager{
		logDir:   logDir,
		maxSize:  10 * 1024 * 1024, // 10MB
		maxFiles: 10,
	}
}

// Entry represents a log entry
type Entry struct {
	Timestamp time.Time
	Level     string
	Package   string
	Message   string
	Command   string
	Output    string
	Error     string
}

// Writer implements io.Writer for log streaming
type Writer struct {
	manager *Manager
	pkgID   string
	file    *os.File
}

// NewWriter creates a new log writer for a package
func (m *Manager) NewWriter(pkgID string) (*Writer, error) {
	if err := os.MkdirAll(m.logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := filepath.Join(m.logDir, fmt.Sprintf("%s-%s.log", pkgID, time.Now().Format("20060102-150405")))

	file, err := os.Create(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return &Writer{
		manager: m,
		pkgID:   pkgID,
		file:    file,
	}, nil
}

// Write implements io.Writer
func (w *Writer) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s", timestamp, string(p))
	return w.file.WriteString(line)
}

// Close closes the log writer
func (w *Writer) Close() error {
	return w.file.Close()
}

// Log writes a log entry
func (m *Manager) Log(entry Entry) error {
	if err := os.MkdirAll(m.logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := filepath.Join(m.logDir, fmt.Sprintf("%s.log", entry.Package))

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)

	if entry.Command != "" {
		logLine += fmt.Sprintf("Command: %s\n", entry.Command)
	}

	if entry.Output != "" {
		logLine += fmt.Sprintf("Output:\n%s\n", entry.Output)
	}

	if entry.Error != "" {
		logLine += fmt.Sprintf("Error: %s\n", entry.Error)
	}

	if _, err := file.WriteString(logLine); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	// Rotate if needed
	m.rotateLog(logFile)

	return nil
}

// GetLogs retrieves logs for a package
func (m *Manager) GetLogs(pkgID string, limit int) ([]Entry, error) {
	logFile := filepath.Join(m.logDir, fmt.Sprintf("%s.log", pkgID))

	data, err := os.ReadFile(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	// Simple log parsing (can be enhanced)
	_ = data
	_ = limit

	return []Entry{}, nil
}

// ListLogFiles lists all log files
func (m *Manager) ListLogFiles() ([]string, error) {
	entries, err := os.ReadDir(m.logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			files = append(files, filepath.Join(m.logDir, entry.Name()))
		}
	}

	return files, nil
}

// CleanOldLogs removes old log files
func (m *Manager) CleanOldLogs(days int) error {
	entries, err := os.ReadDir(m.logDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -days)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			path := filepath.Join(m.logDir, entry.Name())
			os.Remove(path)
		}
	}

	return nil
}

func (m *Manager) rotateLog(logFile string) {
	info, err := os.Stat(logFile)
	if err != nil {
		return
	}

	if info.Size() < m.maxSize {
		return
	}

	// Rotate by renaming
	timestamp := time.Now().Format("20060102-150405")
	newName := fmt.Sprintf("%s.%s", logFile, timestamp)

	os.Rename(logFile, newName)

	// Clean old rotated files
	m.cleanRotatedFiles(logFile)
}

func (m *Manager) cleanRotatedFiles(baseFile string) {
	pattern := baseFile + ".*"
	matches, _ := filepath.Glob(pattern)

	if len(matches) <= m.maxFiles {
		return
	}

	// Remove oldest files
	for i := 0; i < len(matches)-m.maxFiles; i++ {
		os.Remove(matches[i])
	}
}
