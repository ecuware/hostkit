package checker

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"hostkit/internal/config"
)

// SystemInfo holds system information
type SystemInfo struct {
	OS           string
	Version      string
	Architecture string
	TotalRAM     string
	FreeRAM      string
	TotalDisk    string
	FreeDisk     string
	CPUCount     int
}

// GetSystemInfo retrieves system information
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{
		Architecture: runtime.GOARCH,
	}

	// Detect OS
	if err := detectOS(info); err != nil {
		return nil, fmt.Errorf("failed to detect OS: %w", err)
	}

	// Get memory info
	if err := getMemoryInfo(info); err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	// Get disk info
	if err := getDiskInfo(info); err != nil {
		return nil, fmt.Errorf("failed to get disk info: %w", err)
	}

	// Get CPU count
	info.CPUCount = runtime.NumCPU()

	return info, nil
}

func detectOS(info *SystemInfo) error {
	// Try to detect using /etc/os-release
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	content := string(output)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			info.OS = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			info.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}

	return nil
}

func getMemoryInfo(info *SystemInfo) error {
	cmd := exec.Command("free", "-h")
	output, err := cmd.Output()
	if err != nil {
		// Try alternative method
		return getMemoryInfoFromProc(info)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				info.TotalRAM = fields[1]
			}
			if len(fields) >= 4 {
				info.FreeRAM = fields[3]
			}
			break
		}
	}

	return nil
}

func getMemoryInfoFromProc(info *SystemInfo) error {
	cmd := exec.Command("cat", "/proc/meminfo")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, _ := strconv.Atoi(fields[1])
				info.TotalRAM = fmt.Sprintf("%dMB", kb/1024)
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, _ := strconv.Atoi(fields[1])
				info.FreeRAM = fmt.Sprintf("%dMB", kb/1024)
			}
		}
	}

	return nil
}

func getDiskInfo(info *SystemInfo) error {
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") && !strings.HasPrefix(line, "Filesystem") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				info.TotalDisk = fields[1]
				info.FreeDisk = fields[3]
			}
			break
		}
	}

	return nil
}

// CheckRequirements verifies if system meets package requirements
func CheckRequirements(info *SystemInfo, requirements *config.Requirements) []string {
	var issues []string

	// Check OS compatibility
	if len(requirements.OS) > 0 {
		supported := false
		for osName, versions := range requirements.OS {
			if info.OS == osName {
				for _, ver := range versions {
					if strings.HasPrefix(info.Version, ver) {
						supported = true
						break
					}
				}
			}
		}
		if !supported {
			issues = append(issues, fmt.Sprintf("OS %s %s is not in the supported list", info.OS, info.Version))
		}
	}

	// Check architecture
	if len(requirements.Architecture) > 0 {
		supported := false
		for _, arch := range requirements.Architecture {
			if info.Architecture == arch {
				supported = true
				break
			}
		}
		if !supported {
			issues = append(issues, fmt.Sprintf("Architecture %s is not supported", info.Architecture))
		}
	}

	return issues
}
