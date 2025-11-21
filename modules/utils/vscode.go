package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallVSCode() error {
	if isVSCodeInstalled() {
		fmt.Printf("[+] VS Code đã cài đặt... bỏ qua\n")
		return nil
	}

	fmt.Printf("[*] Đang cài đặt VS Code...\n")

	vscodeUrl := "https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-user"
	vscodeInstaller := filepath.Join(os.Getenv("TEMP"), "VSCodeUserSetup.exe")

	downloadFile(vscodeUrl, vscodeInstaller)

	vscodeArgs := []string{
		"/VERYSILENT",
		"/SP-",
		"/MERGETASKS=!runcode,addcontextmenufiles,addcontextmenufolders,associatewithfiles,addtopath",
	}

	cmd := exec.Command(vscodeInstaller, vscodeArgs...)
	cmd.Run()

	os.Remove(vscodeInstaller)

	return nil
}

func isVSCodeInstalled() bool {
	if _, err := exec.LookPath("code"); err == nil {
		return true
	}

	vscodePath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Microsoft VS Code", "Code.exe")
	if _, err := os.Stat(vscodePath); err == nil {
		return true
	}

	return false
}
