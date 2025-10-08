package chocolatey

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type packageInfo struct {
	name string
}

func (p packageInfo) install() error {
	if isPackageInstalled(p.name) {
		fmt.Printf("[+] %s đã cài đặt... bỏ qua\n", p.name)
		return nil
	}
	fmt.Printf("[*] Đang cài đặt %s...\n", p.name)
	return installPackage(p.name)
}

var packages = []packageInfo{
	{"eza"},
	{"mingw"},
	{"make"},
	{"temurin"},
	{"oh-my-posh"},
}

func updatePath() {
	fmt.Printf("[*] Đang cập nhật PATH...\n")
	exec.Command("powershell", "-Command", "$env:Path = [System.Environment]::GetEnvironmentVariable('Path', 'Machine') + ';' + [System.Environment]::GetEnvironmentVariable('Path', 'User')").Run()
}

func refreshEnvironment() {
	exec.Command("cmd", "/c", "refreshenv").Run()
}

func updateProcessEnvironment() {
	systemPath := exec.Command("powershell", "-Command", "[System.Environment]::GetEnvironmentVariable('PATH', 'Machine')")
	if output, err := systemPath.Output(); err == nil {
		systemPathStr := strings.TrimSpace(string(output))
		userPath := exec.Command("powershell", "-Command", "[System.Environment]::GetEnvironmentVariable('PATH', 'User')")
		if userOutput, err := userPath.Output(); err == nil {
			userPathStr := strings.TrimSpace(string(userOutput))
			newPath := systemPathStr + ";" + userPathStr
			os.Setenv("PATH", newPath)
		}
	}
}

func runPowerShellCommand(command string) {
	exec.Command("powershell", "-Command", command).Run()
}

func installIfMissing(packageName string, installFunc func() error) {
	if err := installFunc(); err != nil {
		fmt.Printf("[!] Cài đặt %s thất bại: %v\n", packageName, err)
	}
}

func InstallGit(gitName, gitEmail string) error {
	if isGitInstalled() {
		fmt.Printf("[+] Phát hiện Git, đang gỡ phiên bản cũ...\n")
		uninstallGit()
	}

	fmt.Printf("[*] Đang cài đặt Git...\n")
	if err := installPackageWithParams("git.install", "'/GitAndUnixToolsOnPath /NoShellIntegration /NoGuiHereIntegration'"); err != nil {
		return fmt.Errorf("cài đặt git thất bại: %v", err)
	}

	updatePath()
	refreshEnvironment()
	updateProcessEnvironment()
	time.Sleep(3 * time.Second)

	return configureGit(gitName, gitEmail)
}

func configureGit(gitName, gitEmail string) error {
	if gitName == "" || gitEmail == "" {
		fmt.Printf("[!] Bỏ qua cấu hình Git - thông tin không đầy đủ\n")
		return nil
	}

	fmt.Printf("[*] Đang cấu hình Git...\n")

	configs := map[string]string{
		"user.name":            gitName,
		"user.email":           gitEmail,
		"init.defaultBranch":   "main",
		"push.autoSetupRemote": "true",
	}

	for key, value := range configs {
		if err := runGitCommand("config", "--global", key, value); err != nil {
			return fmt.Errorf("cấu hình git %s thất bại: %v", key, err)
		}
	}
	return nil
}

func runGitCommand(args ...string) error {
	gitPath := "C:\\Program Files\\Git\\cmd\\git.exe"
	return exec.Command(gitPath, args...).Run()
}

func isGitInstalled() bool {
	if _, err := exec.LookPath("git"); err == nil {
		return true
	}
	if _, err := os.Stat("C:\\Program Files\\Git\\cmd\\git.exe"); err == nil {
		return true
	}
	return false
}

func uninstallGit() {
	exec.Command("C:\\Program Files\\Git\\unins000.exe", "/SILENT").Run()
}

func isPackageInstalled(packageName string) bool {
	switch packageName {
	case "temurin":
		_, err := exec.LookPath("java")
		return err == nil
	case "mingw":
		_, err := exec.LookPath("g++")
		return err == nil
	}
	_, err := exec.LookPath(packageName)
	return err == nil
}

func ensureChocoInstalled() error {
	if !isChocoInstalled() {
		return installChoco()
	}
	return nil
}

func installPackage(packageName string) error {
	if err := ensureChocoInstalled(); err != nil {
		return err
	}
	return exec.Command("choco", "install", packageName, "-y").Run()
}

func installPackageWithParams(packageName, params string) error {
	if err := ensureChocoInstalled(); err != nil {
		return err
	}
	return exec.Command("choco", "install", packageName, "--params", params, "-y", "--force").Run()
}

func isChocoInstalled() bool {
	_, err := exec.LookPath("choco")
	return err == nil
}

func installChoco() error {
	if isChocoInstalled() {
		return nil
	}

	fmt.Printf("[*] Đang cài đặt Chocolatey...\n")

	runPowerShellCommand("if (Test-Path 'C:\\ProgramData\\chocolatey') { Remove-Item -Path 'C:\\ProgramData\\chocolatey' -Recurse -Force }")

	cmd := exec.Command("powershell", "-NoProfile", "-InputFormat", "None", "-ExecutionPolicy", "Bypass", "-Command", "[System.Net.ServicePointManager]::SecurityProtocol = 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cài đặt chocolatey thất bại: %v", err)
	}

	updatePath()
	return nil
}

func InstallDevSetup(gitName, gitEmail string) error {
	if err := InstallGit(gitName, gitEmail); err != nil {
		fmt.Printf("[!] Cài Git thất bại: %v\n", err)
		return err
	}

	for _, pkg := range packages {
		installIfMissing(pkg.name, pkg.install)
	}
	return nil
}
