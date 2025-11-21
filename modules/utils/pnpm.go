package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallPNPM(args ...string) (string, error) {
	pnpmPath := pnpmExecutablePath()
	installed := isPNPMInstalled()

	if !installed && len(args) > 0 {
		return "", fmt.Errorf("pnpm chưa cài, path %s k tồn tại", pnpmPath)
	}

	if !installed {
		fmt.Printf("[*] Đang cài đặt pnpm...\n")
		if err := installPNPM(); err != nil {
			return "", err
		}
		if !isPNPMInstalled() {
			return "", fmt.Errorf("install pnpm xong mà k thấy executable")
		}
	} else {
		fmt.Printf("[+] pnpm đã có sẵn, config lại\n")
	}

	if len(args) > 0 {
		if err := runPNPMCommand(pnpmPath, args...); err != nil {
			return "", err
		}
		return pnpmPath, nil
	}

	if err := configurePNPM(pnpmPath); err != nil {
		return "", err
	}

	return pnpmPath, nil
}

func isPNPMInstalled() bool {
	if _, err := os.Stat(pnpmExecutablePath()); err == nil {
		return true
	}

	if _, err := exec.LookPath("pnpm"); err == nil {
		return true
	}

	return false
}

func installPNPM() error {
	cmd := exec.Command("powershell", "-Command", "Invoke-WebRequest https://get.pnpm.io/install.ps1 -UseBasicParsing | Invoke-Expression")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("install pnpm fail: %s: %w", string(out), err)
	}

	return nil
}

func configurePNPM(pnpmPath string) error {
	if err := runPNPMCommand(pnpmPath, "env", "use", "--global", "jod"); err != nil {
		return err
	}

	if err := runPNPMCommand(pnpmPath, "config", "set", "save-prefix", ""); err != nil {
		return err
	}

	return nil
}

func runPNPMCommand(pnpmPath string, args ...string) error {
	cmd := exec.Command(pnpmPath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pnpm cmd fail: %s: %w", string(out), err)
	}
	return nil
}

func pnpmExecutablePath() string {
	return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "pnpm", "pnpm.CMD")
}
