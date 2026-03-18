package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager handles backup operations
type Manager struct {
	backupDir string
}

// NewManager creates a new backup manager
func NewManager(backupDir string) *Manager {
	if backupDir == "" {
		backupDir = "/var/backups/hostkit"
	}
	return &Manager{backupDir: backupDir}
}

// BackupInfo holds backup metadata
type BackupInfo struct {
	ID          string
	Name        string
	Type        string // config, data, full
	Source      string
	Destination string
	Size        int64
	CreatedAt   time.Time
	Files       int
}

// CreateBackup creates a backup
func (m *Manager) CreateBackup(name, backupType string, sources []string) (*BackupInfo, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupID := fmt.Sprintf("%s-%s-%s", name, backupType, timestamp)
	backupPath := filepath.Join(m.backupDir, backupID+".tar.gz")

	// Create backup directory
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create tar.gz file
	file, err := os.Create(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	info := &BackupInfo{
		ID:          backupID,
		Name:        name,
		Type:        backupType,
		Destination: backupPath,
		CreatedAt:   time.Now(),
	}

	// Add files to archive
	for _, source := range sources {
		if err := m.addPathToArchive(tarWriter, source, &info.Files); err != nil {
			return nil, fmt.Errorf("failed to backup %s: %w", source, err)
		}
	}

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get backup size: %w", err)
	}
	info.Size = stat.Size()

	return info, nil
}

func (m *Manager) addPathToArchive(tw *tar.Writer, path string, fileCount *int) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return m.addDirectoryToArchive(tw, path, fileCount)
	}

	return m.addFileToArchive(tw, path, fileCount)
}

func (m *Manager) addDirectoryToArchive(tw *tar.Writer, dir string, fileCount *int) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return m.addFileToArchive(tw, path, fileCount)
	})
}

func (m *Manager) addFileToArchive(tw *tar.Writer, filePath string, fileCount *int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return err
	}

	header.Name = filePath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	*fileCount++
	return nil
}

// ListBackups lists all backups
func (m *Manager) ListBackups() ([]BackupInfo, error) {
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tar.gz") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			ID:        strings.TrimSuffix(entry.Name(), ".tar.gz"),
			Name:      extractBackupName(entry.Name()),
			Type:      extractBackupType(entry.Name()),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return backups, nil
}

func extractBackupName(filename string) string {
	parts := strings.Split(filename, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

func extractBackupType(filename string) string {
	parts := strings.Split(filename, "-")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}

// RestoreBackup restores a backup
func (m *Manager) RestoreBackup(backupID, targetDir string) error {
	backupPath := filepath.Join(m.backupDir, backupID+".tar.gz")

	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to read gzip: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		targetPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// DeleteBackup deletes a backup
func (m *Manager) DeleteBackup(backupID string) error {
	backupPath := filepath.Join(m.backupDir, backupID+".tar.gz")
	return os.Remove(backupPath)
}

// ScheduleBackup schedules automatic backups (placeholder)
func (m *Manager) ScheduleBackup(name, backupType string, sources []string, interval string) error {
	// This would typically create a cron job
	// For now, just return nil
	return nil
}
