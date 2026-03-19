package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

// Server represents a remote server
type Server struct {
	Name     string   `yaml:"name"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	User     string   `yaml:"user"`
	Tags     []string `yaml:"tags"`
	Password string   `yaml:"password,omitempty"`
}

// Cluster represents a group of servers
type Cluster struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Servers     []Server `yaml:"servers"`
}

// Settings contains global configuration
type Settings struct {
	ParallelLimit int    `yaml:"parallel_limit"`
	Timeout       int    `yaml:"timeout"`
	KeyPath       string `yaml:"key_path"`
}

// Config is the main cluster configuration
type Config struct {
	Clusters []Cluster `yaml:"clusters"`
	Settings Settings  `yaml:"settings"`
}

// LoadConfig reads the cluster configuration from file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "/etc/hostkit/cluster/servers.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if config.Settings.ParallelLimit == 0 {
		config.Settings.ParallelLimit = 5
	}
	if config.Settings.Timeout == 0 {
		config.Settings.Timeout = 30
	}
	if config.Settings.KeyPath == "" {
		config.Settings.KeyPath = "/root/.ssh/cluster_keys/master_key"
	}

	return &config, nil
}

// GetServersByTags returns servers matching the given tags
func (c *Config) GetServersByTags(tags []string) []Server {
	var result []Server
	for _, cluster := range c.Clusters {
		for _, server := range cluster.Servers {
			if matchesTags(server.Tags, tags) {
				result = append(result, server)
			}
		}
	}
	return result
}

// GetAllServers returns all configured servers
func (c *Config) GetAllServers() []Server {
	var result []Server
	for _, cluster := range c.Clusters {
		result = append(result, cluster.Servers...)
	}
	return result
}

// GetServerByName returns a server by name
func (c *Config) GetServerByName(name string) (*Server, error) {
	for _, cluster := range c.Clusters {
		for _, server := range cluster.Servers {
			if server.Name == name {
				return &server, nil
			}
		}
	}
	return nil, fmt.Errorf("server not found: %s", name)
}

func matchesTags(serverTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, filter := range filterTags {
		for _, tag := range serverTags {
			if tag == filter {
				return true
			}
		}
	}
	return false
}

// SSHClient manages SSH connections
type SSHClient struct {
	config  *ssh.ClientConfig
	keyPath string
}

// NewSSHClient creates a new SSH client
func NewSSHClient(keyPath string) (*SSHClient, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	return &SSHClient{
		config: &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		},
		keyPath: keyPath,
	}, nil
}

// Execute runs a command on a remote server
func (s *SSHClient) Execute(server Server, command string) (string, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), s.config)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// ExecuteParallel runs commands on multiple servers in parallel
func ExecuteParallel(servers []Server, command string, keyPath string, parallelLimit int) map[string]string {
	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, parallelLimit)

	client, err := NewSSHClient(keyPath)
	if err != nil {
		for _, server := range servers {
			results[server.Name] = fmt.Sprintf("ERROR: %v", err)
		}
		return results
	}

	for _, server := range servers {
		wg.Add(1)
		go func(srv Server) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			output, err := client.Execute(srv, command)
			mu.Lock()
			if err != nil {
				results[srv.Name] = fmt.Sprintf("ERROR: %v\n%s", err, output)
			} else {
				results[srv.Name] = output
			}
			mu.Unlock()
		}(server)
	}

	wg.Wait()
	return results
}

// Ping checks if servers are reachable
func PingServers(servers []Server, keyPath string) map[string]bool {
	results := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	client, err := NewSSHClient(keyPath)
	if err != nil {
		for _, server := range servers {
			results[server.Name] = false
		}
		return results
	}

	for _, server := range servers {
		wg.Add(1)
		go func(srv Server) {
			defer wg.Done()
			_, err := client.Execute(srv, "echo 'pong'")
			mu.Lock()
			results[srv.Name] = (err == nil)
			mu.Unlock()
		}(server)
	}

	wg.Wait()
	return results
}

// SaveConfig writes configuration to file
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
