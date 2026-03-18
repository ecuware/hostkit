package detector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"hostkit/internal/config"
)

// VersionDetector interface for version detection strategies
type VersionDetector interface {
	Detect(ctx context.Context, source config.VersionSource) (string, error)
}

// GitHubReleaseDetector detects version from GitHub releases
type GitHubReleaseDetector struct {
	client *http.Client
}

// NewGitHubReleaseDetector creates a new GitHub detector
func NewGitHubReleaseDetector() *GitHubReleaseDetector {
	return &GitHubReleaseDetector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GitHubRelease represents GitHub release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// Detect fetches latest version from GitHub releases
func (d *GitHubReleaseDetector) Detect(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release: %w", err)
	}

	// Clean version string (remove 'v' prefix if present)
	version := regexp.MustCompile(`^v`).ReplaceAllString(release.TagName, "")
	return version, nil
}

// URLScrapeDetector scrapes version from a URL
type URLScrapeDetector struct {
	client *http.Client
}

// NewURLScrapeDetector creates a new URL scraper
func NewURLScrapeDetector() *URLScrapeDetector {
	return &URLScrapeDetector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Detect scrapes version from URL using regex
func (d *URLScrapeDetector) Detect(ctx context.Context, url, pattern string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %w", err)
	}

	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("version pattern not found in page")
	}

	return matches[1], nil
}

// DetectorManager manages multiple detection strategies
type DetectorManager struct {
	githubDetector *GitHubReleaseDetector
	scrapeDetector *URLScrapeDetector
}

// NewDetectorManager creates a new detector manager
func NewDetectorManager() *DetectorManager {
	return &DetectorManager{
		githubDetector: NewGitHubReleaseDetector(),
		scrapeDetector: NewURLScrapeDetector(),
	}
}

// DetectVersion detects version using the appropriate strategy
func (m *DetectorManager) DetectVersion(ctx context.Context, source config.VersionSource) (string, error) {
	switch source.Type {
	case "github_release":
		return m.githubDetector.Detect(ctx, source.Owner, source.Repo)
	case "url_scrape":
		return m.scrapeDetector.Detect(ctx, source.URL, source.Regex)
	case "static":
		return "", fmt.Errorf("static version detection requires manual configuration")
	default:
		return "", fmt.Errorf("unknown version source type: %s", source.Type)
	}
}
