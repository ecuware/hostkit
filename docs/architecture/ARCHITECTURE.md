# HostKit Architecture

This document describes the high-level architecture of HostKit, explaining how the different components work together.

## Overview

HostKit is a Go-based CLI tool that provides an interactive TUI (Terminal User Interface) for managing hosting server software. It uses a plugin-based architecture where packages are defined as YAML files and installed through shell scripts.

```
┌─────────────────────────────────────────────────────────────┐
│                        User Layer                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │   CLI    │  │   TUI    │  │  Config  │                  │
│  │ Commands │  │(BubbleTea│  │   YAML   │                  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘                  │
└───────┼─────────────┼─────────────┼─────────────────────────┘
        │             │             │
        └─────────────┼─────────────┘
                      │
┌─────────────────────┼──────────────────────────────────────┐
│                   Core Layer                                │
│  ┌──────────────────┴──────────────────┐                   │
│  │         Installer Engine             │                   │
│  │  ┌──────────┬──────────┬──────────┐ │                   │
│  │  │Resolver  │ Executor │ Logger   │ │                   │
│  │  └──────────┴──────────┴──────────┘ │                   │
│  └─────────────────────────────────────┘                   │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Config    │  │   Monitor    │  │   Cluster    │      │
│  │   Parser     │  │   (Stats)    │  │   Manager    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
        │
        │
┌───────┼────────────────────────────────────────────────────┐
│       │              Package Layer                          │
│       │   ┌─────────────────────────────────────────┐      │
│       └──▶│         Package Registry                │      │
│           │  ┌────────┬────────┬────────┬─────────┐ │      │
│           │  │Panels  │Databases│Servers │Security│ │      │
│           │  └────────┴────────┴────────┴─────────┘ │      │
│           └─────────────────────────────────────────┘      │
│                                                             │
│           ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│           │  Shell   │  │  Script  │  │  Package │        │
│           │Installers│  │  Remote  │  │ Manager  │        │
│           └──────────┘  └──────────┘  └──────────┘        │
└────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. User Interface Layer

#### CLI (`cmd/hostkit/main.go`)
- Entry point for the application
- Uses Cobra framework for command handling
- Commands: `install`, `list`, `tui`, `version`, etc.

#### TUI (`internal/tui/`)
- Built with Bubble Tea (Elm Architecture for Go)
- Components:
  - `main.go` - Main TUI loop
  - `model.go` - State management
  - `package_list.go` - Package browsing
  - `install.go` - Installation progress
  - `monitor.go` - System monitoring view
  - `dialog.go` - User confirmations

**Key Features:**
- Keyboard navigation (↑↓←→, Enter, ESC)
- Real-time updates
- Color-coded status indicators
- Responsive layout

### 2. Installer Engine (`internal/installer/`)

The heart of HostKit - manages package installations.

#### Resolver (`resolver.go`)
- Parses package dependencies
- Detects circular dependencies
- Determines installation order
- Resolves conflicts

```go
type Resolver struct {
    packages map[string]*config.Package
    graph    *DependencyGraph
}

func (r *Resolver) Resolve(packageID string) ([]string, error)
```

#### Executor (`executor.go`)
- Executes installation scripts
- Manages process lifecycle
- Captures stdout/stderr
- Handles timeouts

```go
type Executor struct {
    timeout time.Duration
    logger  *Logger
}

func (e *Executor) Run(script string, env map[string]string) error
```

#### Logger (`logger.go`)
- Real-time log output
- Progress tracking
- Error capture
- File persistence

### 3. Configuration System (`internal/config/`)

#### Package Definition
```yaml
# configs/nginx.yaml
id: nginx
name: "Nginx"
category: webserver
version:
  current: "1.25.3"
  source:
    type: github_release
    owner: nginx
    repo: nginx
requirements:
  min_ram: "50MB"
  ports: [80, 443]
install:
  method: shell
  script: |
    apt-get update
    apt-get install -y nginx
```

#### Parser (`types.go`)
- YAML unmarshaling
- Validation
- Default value injection
- Template processing

### 4. Monitoring (`internal/monitor/`)

System resource monitoring for the TUI.

```go
type SystemStats struct {
    CPU    CPUStats
    Memory MemoryStats
    Disk   DiskStats
    Network NetworkStats
    Uptime time.Duration
}
```

**Data Sources:**
- `/proc/stat` - CPU
- `/proc/meminfo` - Memory
- `/proc/diskstats` - Disk I/O
- `/proc/net/dev` - Network

### 5. Cluster Management (`internal/cluster/`)

Multi-server SSH orchestration.

```go
type ClusterManager struct {
    config  *ClusterConfig
    sshKey  string
    clients map[string]*SSHClient
}

func (c *ClusterManager) ExecuteParallel(
    servers []Server, 
    command string,
) map[string]Result
```

## Package System

### Package Categories

Each category is a subdirectory in `configs/`:

```
configs/
├── panels/          # Hosting control panels
│   ├── aapanel.yaml
│   ├── cyberpanel.yaml
│   └── ...
├── databases/       # Database servers
│   ├── mariadb.yaml
│   ├── postgresql.yaml
│   └── ...
├── webservers/      # Web servers
│   ├── nginx.yaml
│   └── ...
├── security/        # Security tools
├── services/        # Additional services
├── monitoring/      # Monitoring tools
├── vpn/            # VPN solutions
└── cluster/        # Cluster management
```

### Installation Methods

1. **Shell Scripts** - Inline bash scripts
2. **Remote Scripts** - Download and execute
3. **Package Managers** - apt, yum, etc.
4. **Docker** - Container-based deployment

### Version Detection

Packages can auto-detect latest versions:

```yaml
version:
  source:
    type: github_release  # Check GitHub releases
    owner: nginx
    repo: nginx
```

Supported types:
- `github_release` - GitHub API
- `api` - Generic REST API
- `url_scrape` - Regex from webpage
- `static` - Manual version

## Data Flow

### Installation Flow

```
1. User selects package
   ↓
2. TUI displays details
   ↓
3. Resolver builds dependency tree
   ↓
4. User confirms
   ↓
5. Executor runs pre-install checks
   ↓
6. Installation script executes
   ↓
7. Post-install hooks run
   ↓
8. Status updated
```

### Configuration Loading

```
1. HostKit starts
   ↓
2. Scan configs/ directory
   ↓
3. Parse each YAML file
   ↓
4. Validate structure
   ↓
5. Build package registry
   ↓
6. Index by category
```

## Error Handling

### Levels

1. **Fatal** - Cannot continue (abort)
2. **Error** - Operation failed (rollback)
3. **Warning** - Issue encountered (continue)
4. **Info** - Status update

### Recovery

```go
// Transaction-like installation
func InstallWithRollback(pkg Package) error {
    // 1. Create backup point
    backup := createBackup()
    
    // 2. Attempt installation
    if err := install(pkg); err != nil {
        // 3. Rollback on failure
        backup.restore()
        return err
    }
    
    // 4. Clean up backup on success
    backup.delete()
    return nil
}
```

## Security Considerations

### Script Execution
- All scripts run as root (required for system packages)
- Scripts are validated before execution
- No arbitrary code execution from untrusted sources

### SSH Cluster
- Key-based authentication only
- No password storage
- Connections time out
- Keys generated per-cluster

### File Permissions
- Config files: 644 (readable by all, writable by owner)
- SSH keys: 600 (owner only)
- Scripts: 755 (executable)

## Performance

### Optimizations

1. **Parallel Execution** - Multiple independent packages
2. **Caching** - Version checks cached
3. **Lazy Loading** - Configs loaded on demand
4. **Streaming** - Real-time output

### Resource Usage

- **Binary Size:** ~10MB
- **Memory:** ~50MB idle
- **CPU:** Minimal (<1% when idle)
- **Disk:** ~100MB (configs + logs)

## Extension Points

### Adding New Categories

1. Create directory: `configs/newcategory/`
2. Add YAML files
3. Update category list
4. No code changes needed!

### Custom Installers

```yaml
install:
  method: custom
  plugin: my_custom_installer
  params:
    key: value
```

### Hooks

```yaml
hooks:
  pre_install: |
    echo "Before installation"
  post_install: |
    echo "After installation"
  on_error: |
    echo "Installation failed"
```

## Testing Strategy

### Unit Tests
- Individual functions
- Mock external dependencies
- Fast execution

### Integration Tests
- End-to-end package installation
- Real system interaction
- CI/CD pipeline

### Manual Testing
- Different OS versions
- Various hardware specs
- Edge cases

## Future Architecture

### Planned Improvements

1. **Plugin System** - Third-party extensions
2. **Web API** - RESTful API for remote management
3. **State Management** - Database-backed state
4. **Rollback** - Better recovery mechanisms
5. **Dry Run** - Test without installing

### Migration Path

Current → Target:
- YAML configs remain compatible
- Gradual feature addition
- Backward compatibility maintained

## Contributing

To modify the architecture:

1. Discuss in an issue first
2. Update this document
3. Maintain backward compatibility
4. Add tests
5. Update examples

See [CONTRIBUTING.md](../CONTRIBUTING.md) for details.
