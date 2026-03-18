package installer

import (
	"fmt"
	"sort"

	"hostkit/internal/config"
)

// DependencyResolver handles dependency resolution
type DependencyResolver struct {
	configs map[string]*config.Config
}

// NewDependencyResolver creates a new resolver
func NewDependencyResolver(configs map[string]*config.Config) *DependencyResolver {
	return &DependencyResolver{
		configs: configs,
	}
}

// Resolve resolves all dependencies for a package
func (r *DependencyResolver) Resolve(pkgID string) ([]string, error) {
	visited := make(map[string]bool)
	result := []string{}

	if err := r.resolveRecursive(pkgID, visited, &result); err != nil {
		return nil, err
	}

	// Reverse to get correct installation order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

func (r *DependencyResolver) resolveRecursive(pkgID string, visited map[string]bool, result *[]string) error {
	if visited[pkgID] {
		return nil
	}

	cfg, exists := r.configs[pkgID]
	if !exists {
		return fmt.Errorf("unknown package: %s", pkgID)
	}

	visited[pkgID] = true

	// Resolve required dependencies first
	for _, dep := range cfg.Dependencies.Required {
		if err := r.resolveRecursive(dep, visited, result); err != nil {
			return fmt.Errorf("dependency %s of %s: %w", dep, pkgID, err)
		}
	}

	*result = append(*result, pkgID)
	return nil
}

// CheckCircular checks for circular dependencies
func (r *DependencyResolver) CheckCircular(pkgID string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	return r.checkCircularRecursive(pkgID, visited, recStack)
}

func (r *DependencyResolver) checkCircularRecursive(pkgID string, visited, recStack map[string]bool) error {
	if recStack[pkgID] {
		return fmt.Errorf("circular dependency detected involving %s", pkgID)
	}

	if visited[pkgID] {
		return nil
	}

	cfg, exists := r.configs[pkgID]
	if !exists {
		return nil // Unknown packages are handled elsewhere
	}

	visited[pkgID] = true
	recStack[pkgID] = true

	for _, dep := range cfg.Dependencies.Required {
		if err := r.checkCircularRecursive(dep, visited, recStack); err != nil {
			return err
		}
	}

	recStack[pkgID] = false
	return nil
}

// GetInstallOrder returns the installation order for multiple packages
func (r *DependencyResolver) GetInstallOrder(pkgIDs []string) ([]string, error) {
	allDeps := make(map[string]bool)

	for _, pkgID := range pkgIDs {
		resolved, err := r.Resolve(pkgID)
		if err != nil {
			return nil, err
		}
		for _, dep := range resolved {
			allDeps[dep] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(allDeps))
	for dep := range allDeps {
		result = append(result, dep)
	}

	// Sort to ensure consistent ordering
	sort.Strings(result)

	return result, nil
}

// IsDependency checks if a package is a dependency of another
func (r *DependencyResolver) IsDependency(pkgID, potentialDep string) bool {
	deps, err := r.Resolve(pkgID)
	if err != nil {
		return false
	}

	for _, dep := range deps {
		if dep == potentialDep {
			return true
		}
	}

	return false
}
