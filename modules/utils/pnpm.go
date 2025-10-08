package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallPNPM() error {
	if isPNPMInstalled() {
		fmt.Printf("[+] pnpm đã cài đặt... bỏ qua\n")
		return nil
	}

	fmt.Printf("[*] Đang cài đặt pnpm...\n")
	installPNPM()
	return nil
}

func isPNPMInstalled() bool {
	pnpmExePath := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "pnpm", "pnpm.EXE")
	if _, err := os.Stat(pnpmExePath); err == nil {
		return true
	}

	pnpmCmdPath := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "pnpm", "pnpm.CMD")
	if _, err := os.Stat(pnpmCmdPath); err == nil {
		return true
	}

	if _, err := exec.LookPath("pnpm"); err == nil {
		return true
	}

	return false
}

func installPNPM() error {
	cmd := exec.Command("powershell", "-Command", "Invoke-WebRequest https://get.pnpm.io/install.ps1 -UseBasicParsing | Invoke-Expression")
	cmd.CombinedOutput()

	cmd = exec.Command("pnpm", "env", "use", "--global", "jod")
	cmd.CombinedOutput()

	cmd = exec.Command("pnpm", "config", "set", "save-prefix", "")
	cmd.CombinedOutput()

	return nil
}
