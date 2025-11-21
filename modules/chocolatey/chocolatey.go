package chocolatey

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	chocoExecutablePath      = "C:\\ProgramData\\chocolatey\\bin\\choco.exe"
	refreshEnvExecutablePath = "C:\\ProgramData\\chocolatey\\bin\\RefreshEnv.cmd"
	chocoBinDir              = "C:\\ProgramData\\chocolatey\\bin"
)

func updatePath() {
	fmt.Printf("[*] Đang cập nhật PATH...\n")
	exec.Command("powershell", "-Command", "$env:Path = [System.Environment]::GetEnvironmentVariable('Path', 'Machine') + ';' + [System.Environment]::GetEnvironmentVariable('Path', 'User')").Run()
}

func refreshEnvironment() {
	if cmd := refreshEnvExecutable(); cmd != "" {
		exec.Command(cmd).Run()
	}
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

func runChoco(args ...string) error {
	cmd := exec.Command(chocoExecutable(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	return runChoco("install", packageName, "-y", "--force")
}

func installPackageWithParams(packageName, params string) error {
	if err := ensureChocoInstalled(); err != nil {
		return err
	}
	return runChoco("install", packageName, "--params", params, "-y", "--force")
}

func isChocoInstalled() bool {
	if _, err := exec.LookPath("choco"); err == nil {
		return true
	}
	if _, err := os.Stat(chocoExecutablePath); err == nil {
		return true
	}
	return false
}

func installChoco() error {
	if isChocoInstalled() {
		return nil
	}

	fmt.Printf("[*] Đang cài đặt Chocolatey...\n")

	runPowerShellCommand("if (Test-Path 'C:\\ProgramData\\chocolatey') { Remove-Item -Path 'C:\\ProgramData\\chocolatey' -Recurse -Force }")

	psScript := `
		Set-ExecutionPolicy Bypass -Scope Process -Force;
		[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072;
		$ProgressPreference = 'SilentlyContinue';
		try {
			Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1')) | Out-Null;
		} finally {
			$ProgressPreference = 'Continue';
		}
	`

	cmd := exec.Command("powershell", "-NoProfile", "-InputFormat", "None", "-ExecutionPolicy", "Bypass", "-Command", psScript)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cài đặt chocolatey thất bại: %v", err)
	}

	updatePath()
	ensureChocoBinInPath()
	refreshEnvironment()
	updateProcessEnvironment()
	return nil
}

func chocoExecutable() string {
	if _, err := os.Stat(chocoExecutablePath); err == nil {
		return chocoExecutablePath
	}
	return "choco"
}

func refreshEnvExecutable() string {
	if _, err := os.Stat(refreshEnvExecutablePath); err == nil {
		return refreshEnvExecutablePath
	}
	return ""
}

func ensureChocoBinInPath() {
	currentPath := os.Getenv("PATH")
	if strings.Contains(strings.ToLower(currentPath), strings.ToLower(chocoBinDir)) {
		return
	}
	if currentPath == "" {
		os.Setenv("PATH", chocoBinDir)
		return
	}
	os.Setenv("PATH", fmt.Sprintf("%s;%s", currentPath, chocoBinDir))
}

func InstallEza() error {
	if isPackageInstalled("eza") {
		fmt.Printf("[+] eza đã cài, bỏ qua\n")
		return nil
	}
	fmt.Printf("[*] Đang cài eza...\n")
	if err := installPackage("eza"); err != nil {
		return fmt.Errorf("cài eza thất bại: %w", err)
	}
	fmt.Printf("[+] Đã cài eza\n")
	return nil
}

func InstallMingw() error {
	if isPackageInstalled("mingw") {
		fmt.Printf("[+] mingw đã cài, bỏ qua\n")
		return nil
	}
	fmt.Printf("[*] Đang cài mingw...\n")
	if err := installPackage("mingw"); err != nil {
		return fmt.Errorf("cài mingw thất bại: %w", err)
	}
	fmt.Printf("[+] Đã cài mingw\n")
	return nil
}

func InstallMake() error {
	if isPackageInstalled("make") {
		fmt.Printf("[+] make đã cài, bỏ qua\n")
		return nil
	}
	fmt.Printf("[*] Đang cài make...\n")
	if err := installPackage("make"); err != nil {
		return fmt.Errorf("cài make thất bại: %w", err)
	}
	fmt.Printf("[+] Đã cài make\n")
	return nil
}

func InstallTemurin() error {
	if isPackageInstalled("temurin") {
		fmt.Printf("[+] temurin đã cài, bỏ qua\n")
		return nil
	}
	fmt.Printf("[*] Đang cài temurin...\n")
	if err := installPackage("temurin"); err != nil {
		return fmt.Errorf("cài temurin thất bại: %w", err)
	}
	fmt.Printf("[+] Đã cài temurin\n")
	return nil
}

func InstallOhMyPosh() error {
	if isPackageInstalled("oh-my-posh") {
		fmt.Printf("[+] oh-my-posh đã cài, bỏ qua\n")
		return nil
	}
	fmt.Printf("[*] Đang cài oh-my-posh...\n")
	if err := installPackage("oh-my-posh"); err != nil {
		return fmt.Errorf("cài oh-my-posh thất bại: %w", err)
	}
	fmt.Printf("[+] Đã cài oh-my-posh\n")
	return nil
}

func handleChocoInstall(step string, fn func() error) {
	if err := fn(); err != nil {
		fmt.Printf("[!] %s lỗi: %v\n", step, err)
	}
}

func InstallDevSetup(gitName, gitEmail string) error {
	if err := InstallGit(gitName, gitEmail); err != nil {
		fmt.Printf("[!] Cài Git thất bại: %v\n", err)
		return err
	}

	handleChocoInstall("cài eza", InstallEza)
	handleChocoInstall("cài mingw", InstallMingw)
	handleChocoInstall("cài make", InstallMake)
	handleChocoInstall("cài temurin", InstallTemurin)
	handleChocoInstall("cài oh-my-posh", InstallOhMyPosh)
	return nil
}
