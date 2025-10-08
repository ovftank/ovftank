package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type GitHubRelease struct {
	Assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func Clink() {
	resp, _ := http.Get("https://api.github.com/repos/chrisant996/clink/releases/latest")
	if resp != nil {
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var release GitHubRelease
		json.Unmarshal(body, &release)

		var setupURL string
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, "setup.exe") {
				setupURL = asset.BrowserDownloadURL
				break
			}
		}

		if setupURL != "" {
			tempDir := os.TempDir()
			installerPath := filepath.Join(tempDir, "clink_setup.exe")

			downloadFile(setupURL, installerPath)
			exec.Command(installerPath, "/S").Run()
			os.Remove(installerPath)

			createClinkAliases()
			configureClinkAutorun()
		}
	}
}

func InstallClink() error {
	clinkPath := filepath.Join(os.Getenv("ProgramFiles(x86)"), "clink", "clink.exe")
	if _, err := os.Stat(clinkPath); err == nil {
		fmt.Printf("[+] clink đã cài đặt... bỏ qua\n")
		return nil
	}

	fmt.Printf("[*] Đang cài đặt clink...\n")
	Clink()
	return nil
}

func downloadFile(url, filepath string) {
	resp, _ := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()

		out, _ := os.Create(filepath)
		if out != nil {
			defer out.Close()
			io.Copy(out, resp.Body)
		}
	}
}

func createClinkAliases() {
	aliasesDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "clink")
	aliasesFile := filepath.Join(aliasesDir, "aliases")

	os.MkdirAll(aliasesDir, 0755)

	aliasesContent := "ls=eza --icons $*\n"

	if _, err := os.Stat(aliasesFile); os.IsNotExist(err) {
		os.WriteFile(aliasesFile, []byte(aliasesContent), 0644)
	}
}

func configureClinkAutorun() {
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	if programFilesX86 == "" {
		programFilesX86 = os.Getenv("ProgramFiles")
	}

	clinkPath := filepath.Join(programFilesX86, "clink", "clink.bat")
	autorunValue := fmt.Sprintf("doskey /macrofile=%%LOCALAPPDATA%%\\clink\\aliases&\"%s\" inject --autorun --quiet", clinkPath)

	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Command Processor`, registry.SET_VALUE)
	if err != nil {
		key, _, err = registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Command Processor`, registry.SET_VALUE)
		if err == nil {
			key.Close()
		}
		return
	}
	defer key.Close()
	key.SetStringValue("Autorun", autorunValue)
}

func OhMyPosh() {
	cmd1 := exec.Command("cmd")
	cmd1.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: `cmd /c "C:\Program Files (x86)\clink\clink.bat" config prompt use oh-my-posh`,
	}
	cmd1.Run()

	clinkExe := `C:\Program Files (x86)\clink\clink.bat`
	themeFile := `C:\Program Files (x86)\oh-my-posh\themes\dracula.omp.json`
	cmd2 := exec.Command("powershell", "-Command",
		fmt.Sprintf("Start-Process -FilePath '%s' -ArgumentList 'set', 'ohmyposh.theme', '%s' -Wait", clinkExe, themeFile))
	cmd2.Run()
}

func ConfigureOhMyPosh() error {
	fmt.Printf("[*] Đang cấu hình Oh My Posh...\n")
	OhMyPosh()
	return nil
}
