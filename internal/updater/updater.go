package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/tui"
	"github.com/charmbracelet/lipgloss"
)

const (
	RepoOwner     = "Mr-Meshky"
	RepoName      = "vify-cli"
	CheckInterval = 24 * time.Hour
)

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type ReleaseInfo struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	HTMLURL     string         `json:"html_url"`
	PublishedAt string         `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

type UpdateCache struct {
	LastChecked   time.Time    `json:"last_checked"`
	LatestRelease *ReleaseInfo `json:"latest_release"`
}

func getCacheFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(home, ".vify")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "update_cache.json")
}

func loadCache() *UpdateCache {
	path := getCacheFilePath()
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cache UpdateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil
	}
	return &cache
}

func saveCache(cache *UpdateCache) {
	path := getCacheFilePath()
	if path == "" {
		return
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}

// FetchLatestRelease queries GitHub API for the latest release
func FetchLatestRelease(timeout time.Duration) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "vify-cli-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// CheckUpdate checks whether a newer version is available.
// If force is false, it uses the cached response if within CheckInterval.
func CheckUpdate(currentVersion string, force bool) (*ReleaseInfo, bool, error) {
	if !force {
		cache := loadCache()
		if cache != nil && cache.LatestRelease != nil && time.Since(cache.LastChecked) < CheckInterval {
			isNewer := CompareVersions(cache.LatestRelease.TagName, currentVersion) > 0
			return cache.LatestRelease, isNewer, nil
		}
	}

	release, err := FetchLatestRelease(5 * time.Second)
	if err != nil {
		// Fallback to cache if available
		if cache := loadCache(); cache != nil && cache.LatestRelease != nil {
			isNewer := CompareVersions(cache.LatestRelease.TagName, currentVersion) > 0
			return cache.LatestRelease, isNewer, nil
		}
		return nil, false, err
	}

	saveCache(&UpdateCache{
		LastChecked:   time.Now(),
		LatestRelease: release,
	})

	isNewer := CompareVersions(release.TagName, currentVersion) > 0
	return release, isNewer, nil
}

// CompareVersions compares two semver strings (e.g. "v1.2.0" and "1.1.5").
// Returns 1 if v1 > v2, -1 if v1 < v2, and 0 if equal.
func CompareVersions(v1, v2 string) int {
	clean := func(v string) []int {
		v = strings.TrimSpace(v)
		v = strings.TrimPrefix(v, "v")
		v = strings.TrimPrefix(v, "V")
		if idx := strings.IndexAny(v, "-+"); idx != -1 {
			v = v[:idx]
		}
		parts := strings.Split(v, ".")
		res := make([]int, 3)
		for i := 0; i < len(parts) && i < 3; i++ {
			n, _ := strconv.Atoi(parts[i])
			res[i] = n
		}
		return res
	}

	parts1 := clean(v1)
	parts2 := clean(v2)

	for i := 0; i < 3; i++ {
		if parts1[i] > parts2[i] {
			return 1
		}
		if parts1[i] < parts2[i] {
			return -1
		}
	}
	return 0
}

// RenderUpdateNotification formats a stylish CLI box notifying about an update
func RenderUpdateNotification(currentVersion, latestVersion string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorWarning).
		Padding(0, 1).
		Margin(1, 0)

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorWarning).Render("⚡ Update Available!")
	versions := fmt.Sprintf("%s → %s",
		lipgloss.NewStyle().Foreground(tui.ColorMuted).Render("v"+strings.TrimPrefix(currentVersion, "v")),
		lipgloss.NewStyle().Foreground(tui.ColorSuccess).Bold(true).Render(latestVersion),
	)
	cmdHint := fmt.Sprintf("Run %s to upgrade.",
		lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true).Render("vify update"),
	)

	content := fmt.Sprintf("%s  %s\n%s", title, versions, cmdHint)
	return boxStyle.Render(content)
}

// FindMatchingAsset finds the asset matching current OS and Architecture
func FindMatchingAsset(release *ReleaseInfo) (*ReleaseAsset, error) {
	osName := runtime.GOOS
	archName := runtime.GOARCH

	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, osName) && strings.Contains(name, archName) {
			if strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip") {
				return &asset, nil
			}
		}
	}

	// Fallback check: OS match
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, osName) && (strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip")) {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("no matching release asset found for %s/%s", osName, archName)
}

// SelfUpdate downloads and installs the latest version into the current binary location
func SelfUpdate(currentVersion string, progressFn func(string)) (*ReleaseInfo, error) {
	if progressFn == nil {
		progressFn = func(string) {}
	}

	progressFn("Checking for latest release...")
	release, isNewer, err := CheckUpdate(currentVersion, true)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}

	if !isNewer && CompareVersions(release.TagName, currentVersion) == 0 {
		return release, nil // Already up to date
	}

	asset, err := FindMatchingAsset(release)
	if err != nil {
		return nil, err
	}

	progressFn(fmt.Sprintf("Downloading %s (%s)...", release.TagName, asset.Name))

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to locate current executable: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve symlink for executable: %w", err)
	}

	// Download archive
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with HTTP %d", resp.StatusCode)
	}

	tempDir, err := os.MkdirTemp("", "vify-update-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, asset.Name)
	outFile, err := os.Create(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to save download archive: %w", err)
	}
	_, err = io.Copy(outFile, resp.Body)
	outFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed writing archive: %w", err)
	}

	progressFn("Extracting binary...")
	extractedBinary, err := extractBinary(archivePath, tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to extract binary: %w", err)
	}

	progressFn("Replacing executable...")
	if err := replaceBinary(extractedBinary, execPath); err != nil {
		return nil, fmt.Errorf("failed to replace executable: %w (try running with sudo)", err)
	}

	return release, nil
}

func extractBinary(archivePath, destDir string) (string, error) {
	binaryName := "vify"
	if runtime.GOOS == "windows" {
		binaryName = "vify.exe"
	}

	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir, binaryName)
	}
	return extractTarGz(archivePath, destDir, binaryName)
}

func extractTarGz(archivePath, destDir, binaryName string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	destBinaryPath := filepath.Join(destDir, binaryName)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		baseName := filepath.Base(header.Name)
		if baseName == binaryName {
			outFile, err := os.OpenFile(destBinaryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()
			return destBinaryPath, nil
		}
	}

	return "", fmt.Errorf("binary '%s' not found inside archive", binaryName)
}

func extractZip(archivePath, destDir, binaryName string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	destBinaryPath := filepath.Join(destDir, binaryName)
	for _, f := range r.File {
		baseName := filepath.Base(f.Name)
		if baseName == binaryName {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			outFile, err := os.OpenFile(destBinaryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				rc.Close()
				return "", err
			}
			_, err = io.Copy(outFile, rc)
			rc.Close()
			outFile.Close()
			if err != nil {
				return "", err
			}
			return destBinaryPath, nil
		}
	}

	return "", fmt.Errorf("binary '%s' not found inside zip archive", binaryName)
}

func replaceBinary(newBinary, targetExec string) error {
	_ = os.Chmod(newBinary, 0755)

	targetDir := filepath.Dir(targetExec)
	tempTarget := filepath.Join(targetDir, fmt.Sprintf(".vify_update_%d", time.Now().UnixNano()))

	src, err := os.Open(newBinary)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(tempTarget, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		_ = os.Remove(tempTarget)
		return err
	}
	dst.Close()

	if runtime.GOOS == "windows" {
		oldBackup := targetExec + ".old"
		_ = os.Remove(oldBackup)
		if err := os.Rename(targetExec, oldBackup); err != nil {
			_ = os.Remove(tempTarget)
			return err
		}
	}

	if err := os.Rename(tempTarget, targetExec); err != nil {
		_ = os.Remove(tempTarget)
		return err
	}

	return nil
}
