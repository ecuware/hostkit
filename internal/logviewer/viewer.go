package logviewer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// LogType represents log type
type LogType string

const (
	LogSystem  LogType = "system"
	LogAuth    LogType = "auth"
	LogNginx   LogType = "nginx"
	LogApache  LogType = "apache"
	LogMySQL   LogType = "mysql"
	LogPackage LogType = "package"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp string
	Level     string
	Service   string
	Message   string
}

// Viewer handles log viewing
type Viewer struct{}

// NewViewer creates a new log viewer
func NewViewer() *Viewer {
	return &Viewer{}
}

// GetAvailableLogs returns available log files
func (v *Viewer) GetAvailableLogs() map[string]string {
	logs := map[string]string{
		"System":   "/var/log/syslog",
		"Auth":     "/var/log/auth.log",
		"Nginx":    "/var/log/nginx/error.log",
		"Apache":   "/var/log/apache2/error.log",
		"MySQL":    "/var/log/mysql/error.log",
		"Fail2ban": "/var/log/fail2ban.log",
		"HostKit":  "/var/log/hostkit/hostkit.log",
	}

	// Check which logs exist
	available := make(map[string]string)
	for name, path := range logs {
		if _, err := os.Stat(path); err == nil {
			available[name] = path
		}
	}

	// Try alternative paths
	alternatives := map[string]string{
		"System": "/var/log/messages",
		"Nginx":  "/var/log/nginx/access.log",
		"MySQL":  "/var/log/mysqld.log",
	}

	for name, path := range alternatives {
		if _, err := os.Stat(path); err == nil {
			if _, exists := available[name]; !exists {
				available[name] = path
			}
		}
	}

	return available
}

// ReadLogs reads log file with filtering
func (v *Viewer) ReadLogs(logPath string, limit int, filter string) ([]LogEntry, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	// Read from end if limit is specified
	if limit > 0 {
		// Use tail command for efficiency
		cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", limit), logPath)
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to read logs: %w", err)
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}

			if filter != "" && !strings.Contains(line, filter) {
				continue
			}

			entry := parseLogLine(line)
			entries = append(entries, entry)
		}
	} else {
		for scanner.Scan() {
			line := scanner.Text()

			if filter != "" && !strings.Contains(line, filter) {
				continue
			}

			entry := parseLogLine(line)
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// FollowLogs follows log file in real-time
func (v *Viewer) FollowLogs(logPath string, filter string, callback func(string)) error {
	cmd := exec.Command("tail", "-f", "-n", "50", logPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if filter != "" && !strings.Contains(line, filter) {
			continue
		}

		callback(line)
	}

	return cmd.Wait()
}

// GetServiceLogs gets logs for a specific service
func (v *Viewer) GetServiceLogs(service string, limit int) ([]LogEntry, error) {
	cmd := exec.Command("journalctl", "-u", service, "-n", fmt.Sprintf("%d", limit), "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get service logs: %w", err)
	}

	var entries []LogEntry
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		entry := parseJournalLine(line, service)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetHostKitLogs gets HostKit specific logs
func (v *Viewer) GetHostKitLogs(limit int) ([]LogEntry, error) {
	logPath := "/var/log/hostkit/install.log"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return []LogEntry{}, nil
	}

	return v.ReadLogs(logPath, limit, "")
}

func parseLogLine(line string) LogEntry {
	entry := LogEntry{Message: line}

	// Try to extract timestamp (common format)
	if len(line) > 15 {
		// Check for date format: Mon DD HH:MM:SS
		if line[3] == ' ' && line[6] == ' ' {
			entry.Timestamp = line[:15]
			line = line[16:]
		}
	}

	// Try to extract log level
	levels := []string{"ERROR", "WARN", "WARNING", "INFO", "DEBUG", "FATAL", "CRITICAL"}
	for _, level := range levels {
		if strings.Contains(line, level) {
			entry.Level = level
			break
		}
	}

	// Extract message
	entry.Message = line

	return entry
}

func parseJournalLine(line, service string) LogEntry {
	entry := LogEntry{
		Service: service,
		Message: line,
	}

	// Journal format: Mon DD HH:MM:SS hostname service[pid]: message
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		entry.Message = strings.TrimSpace(parts[1])

		// Extract timestamp from first part
		fields := strings.Fields(parts[0])
		if len(fields) >= 3 {
			entry.Timestamp = strings.Join(fields[:3], " ")
		}
	}

	// Detect level from message
	if strings.Contains(entry.Message, "error") || strings.Contains(entry.Message, "failed") {
		entry.Level = "ERROR"
	} else if strings.Contains(entry.Message, "warning") || strings.Contains(entry.Message, "warn") {
		entry.Level = "WARN"
	} else {
		entry.Level = "INFO"
	}

	return entry
}

// GetLevelIcon returns icon for log level
func GetLevelIcon(level string) string {
	switch strings.ToUpper(level) {
	case "ERROR", "FATAL", "CRITICAL":
		return "❌"
	case "WARN", "WARNING":
		return "⚠️"
	case "INFO":
		return "ℹ️"
	case "DEBUG":
		return "🐛"
	default:
		return "•"
	}
}
