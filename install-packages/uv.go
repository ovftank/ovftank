package install_packages

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

type UvReceipt struct {
	Binaries       []string          `json:"binaries"`
	BinaryAliases  map[string]string `json:"binary_aliases"`
	Cdylibs        []string          `json:"cdylibs"`
	Cstaticlibs    []string          `json:"cstaticlibs"`
	InstallLayout  string            `json:"install_layout"`
	InstallPrefix  string            `json:"install_prefix"`
	ModifyPath     bool              `json:"modify_path"`
	Provider       ProviderInfo      `json:"provider"`
	Source         SourceInfo        `json:"source"`
	Version        string            `json:"version"`
}

type ProviderInfo struct {
	Source string `json:"source"`
	Version string `json:"version"`
}

type SourceInfo struct {
	AppName     string `json:"app_name"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	ReleaseType string `json:"release_type"`
}

func getUvInstallDir() string {
	// Check for custom install directory
	if installDir := os.Getenv("UV_INSTALL_DIR"); installDir != "" {
		return installDir
	}

	// Check for XDG_BIN_HOME
	if xdgBinHome := os.Getenv("XDG_BIN_HOME"); xdgBinHome != "" {
		return xdgBinHome
	}

	// Check for XDG_DATA_HOME/../bin
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "..", "bin")
	}

	// Default to HOME/.local/bin
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "bin")
}

func getReceiptDir() string {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "uv")
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		return filepath.Join(localAppData, "uv")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "uv")
}

func getPlatformTriple() string {
	arch := runtime.GOARCH
	osName := runtime.GOOS

	if osName != "windows" {
		return ""
	}

	switch arch {
	case "amd64":
		return "x86_64-pc-windows-msvc"
	case "386":
		return "i686-pc-windows-msvc"
	case "arm64":
		return "aarch64-pc-windows-msvc"
	default:
		return "x86_64-pc-windows-msvc" // default fallback
	}
}

func getArtifactName(platform string) string {
	switch platform {
	case "x86_64-pc-windows-msvc":
		return "uv-x86_64-pc-windows-msvc.zip"
	case "i686-pc-windows-msvc":
		return "uv-i686-pc-windows-msvc.zip"
	case "aarch64-pc-windows-msvc":
		return "uv-aarch64-pc-windows-msvc.zip"
	default:
		return "uv-x86_64-pc-windows-msvc.zip"
	}
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func addToPath(installDir string) error {
	// For Windows, we need to modify the registry
	// This is a simplified version - in practice you'd want to use Windows API
	path := os.Getenv("PATH")
	if !strings.Contains(path, installDir) {
		newPath := installDir + ";" + path
		return os.Setenv("PATH", newPath)
	}
	return nil
}

func writeReceipt(receiptDir, installDir, version string) error {
	if err := os.MkdirAll(receiptDir, 0755); err != nil {
		return err
	}

	receipt := UvReceipt{
		Binaries:      []string{"uv.exe", "uvx.exe", "uvw.exe"},
		BinaryAliases: make(map[string]string),
		Cdylibs:       []string{},
		Cstaticlibs:   []string{},
		InstallLayout: "flat",
		InstallPrefix: installDir,
		ModifyPath:    true,
		Provider: ProviderInfo{
			Source:  "cargo-dist",
			Version: "0.30.2",
		},
		Source: SourceInfo{
			AppName:     "uv",
			Name:        "uv",
			Owner:       "astral-sh",
			ReleaseType: "github",
		},
		Version: version,
	}

	receiptFile := filepath.Join(receiptDir, "uv-receipt.json")
	file, err := os.Create(receiptFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(receipt)
}

func InstallUv() {
	if runtime.GOOS != "windows" {
		log.Fatal("UV installer currently only supports Windows")
		return
	}

	// Get latest release info
	resp, err := http.Get("https://api.github.com/repos/astral-sh/uv/releases/latest")
	if err != nil {
		log.Printf("Failed to get latest release: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GitHub API returned status: %d", resp.StatusCode)
		return
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Printf("Failed to decode release info: %v", err)
		return
	}

	version := strings.TrimPrefix(release.TagName, "v")
	if version == "" {
		log.Printf("No version found in release")
		return
	}

	// Determine platform and artifact
	platform := getPlatformTriple()
	artifactName := getArtifactName(platform)

	// Find the correct asset URL
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == artifactName {
			downloadURL = asset.URL
			break
		}
	}

	if downloadURL == "" {
		// Fallback to direct download URL
		downloadURL = fmt.Sprintf("https://github.com/astral-sh/uv/releases/download/%s/%s", release.TagName, artifactName)
	}

	// Get install directory
	installDir := getUvInstallDir()
	if installDir == "" {
		log.Fatal("Could not determine install directory")
		return
	}

	if err := os.MkdirAll(installDir, 0755); err != nil {
		log.Printf("Failed to create install directory: %v", err)
		return
	}

	// Download the archive
	tempFile := filepath.Join(os.TempDir(), artifactName)
	if err := downloadFile(downloadURL, tempFile); err != nil {
		log.Printf("Failed to download UV: %v", err)
		return
	}
	defer os.Remove(tempFile)

	fmt.Printf("Installing UV %s (%s) to %s\n", version, platform, installDir)

	// Extract the zip file
	tempDir, err := os.MkdirTemp("", "uv-install-")
	if err != nil {
		log.Printf("Failed to create temp directory: %v", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Use PowerShell to extract zip (since Go doesn't have built-in zip extraction)
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force", tempFile, tempDir))
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to extract archive: %v", err)
		return
	}

	// Copy binaries to install directory
	binaries := []string{"uv.exe", "uvx.exe", "uvw.exe"}
	for _, binary := range binaries {
		srcPath := filepath.Join(tempDir, binary)
		dstPath := filepath.Join(installDir, binary)

		if _, err := os.Stat(srcPath); err == nil {
			if err := copyFile(srcPath, dstPath); err != nil {
				log.Printf("Failed to copy %s: %v", binary, err)
				continue
			}
			fmt.Printf("  %s\n", binary)
		}
	}

	// Write receipt
	receiptDir := getReceiptDir()
	if err := writeReceipt(receiptDir, installDir, version); err != nil {
		log.Printf("Warning: Failed to write receipt: %v", err)
	}

	// Add to PATH if not disabled
	if os.Getenv("UV_NO_MODIFY_PATH") == "" {
		if err := addToPath(installDir); err != nil {
			log.Printf("Warning: Failed to add to PATH: %v", err)
		} else {
			fmt.Printf("\nUV has been installed to %s\n", installDir)
			fmt.Printf("To add UV to your PATH permanently, restart your shell or run:\n")
			fmt.Printf("    set Path=%s;%%Path%%\n", installDir)
		}
	}

	fmt.Println("UV installation completed successfully!")
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}