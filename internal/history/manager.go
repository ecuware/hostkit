package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single installation history entry
type Entry struct {
	ID          string    `json:"id"`
	PackageID   string    `json:"package_id"`
	PackageName string    `json:"package_name"`
	Version     string    `json:"version"`
	Status      string    `json:"status"` // success, failed, cancelled
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    string    `json:"duration"`
	Error       string    `json:"error,omitempty"`
	LogFile     string    `json:"log_file,omitempty"`
}

// Manager handles installation history
type Manager struct {
	dataDir string
}

// NewManager creates a new history manager
func NewManager(dataDir string) *Manager {
	if dataDir == "" {
		dataDir = "/var/lib/hostkit"
	}
	return &Manager{dataDir: dataDir}
}

// AddEntry adds a new history entry
func (m *Manager) AddEntry(entry Entry) error {
	if err := os.MkdirAll(m.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	historyFile := filepath.Join(m.dataDir, "history.json")

	// Read existing history
	var history []Entry
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}

	// Add new entry
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	history = append(history, entry)

	// Keep only last 100 entries
	if len(history) > 100 {
		history = history[len(history)-100:]
	}

	// Save history
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(historyFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write history: %w", err)
	}

	return nil
}

// GetHistory returns installation history
func (m *Manager) GetHistory(limit int) ([]Entry, error) {
	historyFile := filepath.Join(m.dataDir, "history.json")

	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("failed to read history: %w", err)
	}

	var history []Entry
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal history: %w", err)
	}

	// Reverse order (newest first)
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	if limit > 0 && len(history) > limit {
		history = history[:limit]
	}

	return history, nil
}

// GetPackageHistory returns history for a specific package
func (m *Manager) GetPackageHistory(pkgID string) ([]Entry, error) {
	history, err := m.GetHistory(0)
	if err != nil {
		return nil, err
	}

	var pkgHistory []Entry
	for _, entry := range history {
		if entry.PackageID == pkgID {
			pkgHistory = append(pkgHistory, entry)
		}
	}

	return pkgHistory, nil
}

// GetLastInstall returns the last installation status for a package
func (m *Manager) GetLastInstall(pkgID string) (*Entry, error) {
	history, err := m.GetPackageHistory(pkgID)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, nil
	}

	return &history[0], nil
}

// ClearHistory clears all history
func (m *Manager) ClearHistory() error {
	historyFile := filepath.Join(m.dataDir, "history.json")
	return os.Remove(historyFile)
}
