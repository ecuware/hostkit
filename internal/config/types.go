package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents a service/package configuration
type Config struct {
	ID           string       `yaml:"id"`
	Name         string       `yaml:"name"`
	Category     string       `yaml:"category"`
	Description  string       `yaml:"description"`
	Icon         string       `yaml:"icon"`
	Version      Version      `yaml:"version"`
	Requirements Requirements `yaml:"requirements"`
	Install      Install      `yaml:"install"`
	Dependencies Dependencies `yaml:"dependencies"`
	ConfigFiles  ConfigFiles  `yaml:"config"`
	Uninstall    Uninstall    `yaml:"uninstall"`
	License      License      `yaml:"license"`
	Support      Support      `yaml:"support"`
}

// Version contains version info and detection settings
type Version struct {
	Current string        `yaml:"current"`
	Source  VersionSource `yaml:"source"`
}

// VersionSource defines how to detect latest version
type VersionSource struct {
	Type         string `yaml:"type"` // github_release, url_scrape, static
	Owner        string `yaml:"owner,omitempty"`
	Repo         string `yaml:"repo,omitempty"`
	URL          string `yaml:"url,omitempty"`
	Regex        string `yaml:"regex,omitempty"`
	AssetPattern string `yaml:"asset_pattern,omitempty"`
	FallbackURL  string `yaml:"fallback_url,omitempty"`
	JSONPath     string `yaml:"json_path,omitempty"`
}

// Requirements defines system requirements
type Requirements struct {
	OS            map[string][]string `yaml:"os"`
	MinRAM        string              `yaml:"min_ram,omitempty"`
	MinDisk       string              `yaml:"min_disk,omitempty"`
	Ports         []int               `yaml:"ports,omitempty"`
	Architecture  []string            `yaml:"architecture,omitempty"`
	InstallSize   string              `yaml:"install_size,omitempty"`
	EstimatedTime string              `yaml:"estimated_time,omitempty"`
}

// Install defines installation methods
type Install struct {
	Method       string       `yaml:"method"`
	Script       string       `yaml:"script,omitempty"`
	URL          string       `yaml:"url,omitempty"`
	Args         []string     `yaml:"args,omitempty"`
	Packages     []string     `yaml:"packages,omitempty"`
	Repositories []Repository `yaml:"repositories,omitempty"`
	Ubuntu       *OSInstall   `yaml:"ubuntu,omitempty"`
	CentOS       *OSInstall   `yaml:"centos,omitempty"`
	Debian       *OSInstall   `yaml:"debian,omitempty"`
	PreCheck     []Check      `yaml:"pre_check,omitempty"`
	PostCheck    []Check      `yaml:"post_check,omitempty"`
}

// OSInstall defines OS-specific installation
type OSInstall struct {
	Method       string       `yaml:"method"`
	Packages     []string     `yaml:"packages,omitempty"`
	Repositories []Repository `yaml:"repositories,omitempty"`
	PostInstall  string       `yaml:"post_install,omitempty"`
}

// Repository defines a package repository
type Repository struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Key      string `yaml:"key,omitempty"`
	GPGCheck bool   `yaml:"gpgcheck,omitempty"`
}

// Check defines pre/post installation checks
type Check struct {
	Command  string `yaml:"command,omitempty"`
	Service  string `yaml:"service,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Timeout  int    `yaml:"timeout,omitempty"`
	ErrorMsg string `yaml:"error_msg,omitempty"`
}

// Dependencies lists service dependencies
type Dependencies struct {
	Required []string `yaml:"required,omitempty"`
	Optional []string `yaml:"optional,omitempty"`
}

// ConfigFiles lists configuration files
type ConfigFiles struct {
	Files []ConfigFile `yaml:"files,omitempty"`
}

// ConfigFile represents a config file location
type ConfigFile struct {
	Path        string `yaml:"path"`
	Description string `yaml:"description,omitempty"`
	Template    string `yaml:"template,omitempty"`
}

// Uninstall defines uninstallation method
type Uninstall struct {
	Command string `yaml:"command"`
}

// License defines licensing info
type License struct {
	Required bool   `yaml:"required"`
	Type     string `yaml:"type,omitempty"`
}

// Support defines support links
type Support struct {
	Docs   string `yaml:"docs,omitempty"`
	Forum  string `yaml:"forum,omitempty"`
	Issues string `yaml:"issues,omitempty"`
}

// LoadConfig loads a config from YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// LoadAllConfigs loads all configs from a directory
func LoadAllConfigs(basePath string) (map[string]*Config, error) {
	configs := make(map[string]*Config)

	categories := []string{"panels", "databases", "webservers", "security", "services", "monitoring", "vpn"}

	for _, category := range categories {
		dir := filepath.Join(basePath, category)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Skip if directory doesn't exist
		}

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
				continue
			}

			path := filepath.Join(dir, entry.Name())
			cfg, err := LoadConfig(path)
			if err != nil {
				continue // Skip invalid configs
			}

			configs[cfg.ID] = cfg
		}
	}

	return configs, nil
}

// GetOSInstall returns OS-specific install config
func (c *Config) GetOSInstall(osName string) *OSInstall {
	switch osName {
	case "ubuntu":
		return c.Install.Ubuntu
	case "centos":
		return c.Install.CentOS
	case "debian":
		return c.Install.Debian
	default:
		return nil
	}
}
