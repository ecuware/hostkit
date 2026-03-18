package monitor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// IsLinux checks if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// ErrNotLinux is returned when not running on Linux
var ErrNotLinux = fmt.Errorf("system monitor requires Linux (current OS: %s)", runtime.GOOS)

// SystemStats holds system statistics
type SystemStats struct {
	Timestamp time.Time
	CPU       CPUStats
	Memory    MemoryStats
	Disk      DiskStats
	Network   NetworkStats
	Load      LoadStats
	Uptime    string
}

// CPUStats holds CPU statistics
type CPUStats struct {
	Usage  float64
	User   float64
	System float64
	Idle   float64
	Cores  int
}

// MemoryStats holds memory statistics
type MemoryStats struct {
	Total    uint64
	Used     uint64
	Free     uint64
	Cached   uint64
	Buffers  uint64
	UsagePct float64
}

// DiskStats holds disk statistics
type DiskStats struct {
	Total    uint64
	Used     uint64
	Free     uint64
	UsagePct float64
}

// NetworkStats holds network statistics
type NetworkStats struct {
	Interface string
	RxBytes   uint64
	TxBytes   uint64
	RxSpeed   float64 // bytes per second
	TxSpeed   float64 // bytes per second
}

// LoadStats holds load average
type LoadStats struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// Monitor handles system monitoring
type Monitor struct {
	lastNetStats map[string]NetworkStats
	lastCheck    time.Time
}

// NewMonitor creates a new monitor
func NewMonitor() *Monitor {
	return &Monitor{
		lastNetStats: make(map[string]NetworkStats),
	}
}

// GetStats returns current system stats
func (m *Monitor) GetStats() (*SystemStats, error) {
	if !IsLinux() {
		return nil, ErrNotLinux
	}

	stats := &SystemStats{
		Timestamp: time.Now(),
	}

	// Get CPU stats
	if err := m.getCPUStats(&stats.CPU); err != nil {
		return nil, err
	}

	// Get Memory stats
	if err := m.getMemoryStats(&stats.Memory); err != nil {
		return nil, err
	}

	// Get Disk stats
	if err := m.getDiskStats(&stats.Disk); err != nil {
		return nil, err
	}

	// Get Network stats
	if err := m.getNetworkStats(&stats.Network); err != nil {
		return nil, err
	}

	// Get Load stats
	if err := m.getLoadStats(&stats.Load); err != nil {
		return nil, err
	}

	// Get Uptime
	stats.Uptime = m.getUptime()

	return stats, nil
}

func (m *Monitor) getCPUStats(cpu *CPUStats) error {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}

			user, _ := strconv.ParseFloat(fields[1], 64)
			nice, _ := strconv.ParseFloat(fields[2], 64)
			system, _ := strconv.ParseFloat(fields[3], 64)
			idle, _ := strconv.ParseFloat(fields[4], 64)
			iowait, _ := strconv.ParseFloat(fields[5], 64)
			irq, _ := strconv.ParseFloat(fields[6], 64)
			softirq, _ := strconv.ParseFloat(fields[7], 64)

			total := user + nice + system + idle + iowait + irq + softirq
			cpu.Usage = ((total - idle) / total) * 100
			cpu.User = (user / total) * 100
			cpu.System = (system / total) * 100
			cpu.Idle = (idle / total) * 100
			break
		}
	}

	cpu.Cores = getCPUCores()
	return scanner.Err()
}

func (m *Monitor) getMemoryStats(mem *MemoryStats) error {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024 // Convert from KB to bytes

		switch fields[0] {
		case "MemTotal:":
			mem.Total = value
		case "MemFree:":
			mem.Free = value
		case "Buffers:":
			mem.Buffers = value
		case "Cached:":
			mem.Cached = value
		case "MemAvailable:":
			// Available is more accurate
			mem.Free = value
		}
	}

	mem.Used = mem.Total - mem.Free
	if mem.Total > 0 {
		mem.UsagePct = (float64(mem.Used) / float64(mem.Total)) * 100
	}

	return scanner.Err()
}

func (m *Monitor) getDiskStats(disk *DiskStats) error {
	cmd := exec.Command("df", "-B1", "/")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") && !strings.HasPrefix(line, "Filesystem") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				disk.Total, _ = strconv.ParseUint(fields[1], 10, 64)
				disk.Used, _ = strconv.ParseUint(fields[2], 10, 64)
				disk.Free, _ = strconv.ParseUint(fields[3], 10, 64)

				if disk.Total > 0 {
					disk.UsagePct = (float64(disk.Used) / float64(disk.Total)) * 100
				}
				break
			}
		}
	}

	return nil
}

func (m *Monitor) getNetworkStats(net *NetworkStats) error {
	// Get primary interface
	cmd := exec.Command("ip", "route", "get", "8.8.8.8")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		fields := strings.Fields(lines[0])
		for i, field := range fields {
			if field == "dev" && i+1 < len(fields) {
				net.Interface = fields[i+1]
				break
			}
		}
	}

	// Get interface stats
	if net.Interface != "" {
		rxFile := fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", net.Interface)
		txFile := fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", net.Interface)

		rxData, _ := os.ReadFile(rxFile)
		txData, _ := os.ReadFile(txFile)

		net.RxBytes, _ = strconv.ParseUint(strings.TrimSpace(string(rxData)), 10, 64)
		net.TxBytes, _ = strconv.ParseUint(strings.TrimSpace(string(txData)), 10, 64)

		// Calculate speed
		if lastStats, ok := m.lastNetStats[net.Interface]; ok {
			duration := time.Since(m.lastCheck).Seconds()
			if duration > 0 {
				net.RxSpeed = float64(net.RxBytes-lastStats.RxBytes) / duration
				net.TxSpeed = float64(net.TxBytes-lastStats.TxBytes) / duration
			}
		}

		m.lastNetStats[net.Interface] = *net
	}

	m.lastCheck = time.Now()
	return nil
}

func (m *Monitor) getLoadStats(load *LoadStats) error {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return err
	}

	fields := strings.Fields(string(data))
	if len(fields) >= 3 {
		load.Load1, _ = strconv.ParseFloat(fields[0], 64)
		load.Load5, _ = strconv.ParseFloat(fields[1], 64)
		load.Load15, _ = strconv.ParseFloat(fields[2], 64)
	}

	return nil
}

func (m *Monitor) getUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "Unknown"
	}

	fields := strings.Fields(string(data))
	if len(fields) > 0 {
		uptime, _ := strconv.ParseFloat(fields[0], 64)
		return formatDuration(time.Duration(uptime) * time.Second)
	}

	return "Unknown"
}

func getCPUCores() int {
	cmd := exec.Command("nproc")
	output, err := cmd.Output()
	if err != nil {
		return 1
	}

	cores, _ := strconv.Atoi(strings.TrimSpace(string(output)))
	if cores == 0 {
		return 1
	}
	return cores
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// FormatBytes formats bytes to human readable
func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FormatSpeed formats speed to human readable
func FormatSpeed(bytesPerSec float64) string {
	return FormatBytes(uint64(bytesPerSec)) + "/s"
}
