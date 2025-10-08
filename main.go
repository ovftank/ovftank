package main

import (
	"fmt"
	"os"
	"os/exec"

	"windows-dev-setup/modules/chocolatey"
	"windows-dev-setup/modules/registry"
	"windows-dev-setup/modules/service"
	"windows-dev-setup/modules/utils"
)

func main() {
	exec.Command("cmd", "/c", "color a").Run()
	clearScreen()

	fmt.Printf("\n%s====================================%s\n", "\x1b[32m", "\x1b[0m")
	fmt.Printf("%s>>> WINDOWS DEV SETUP <<<%s\n", "\x1b[32m", "\x1b[0m")
	fmt.Printf("%s====================================%s\n\n", "\x1b[32m", "\x1b[0m")

	var gitName, gitEmail string
	fmt.Printf("%s[INPUT] Nhập username Git: %s", "\x1b[32m", "\x1b[0m")
	fmt.Scanln(&gitName)
	fmt.Printf("%s[INPUT] Nhập email Git: %s", "\x1b[32m", "\x1b[0m")
	fmt.Scanln(&gitEmail)

	if err := chocolatey.InstallDevSetup(gitName, gitEmail); err != nil {
		fmt.Printf("%s[!] Lỗi cài đặt: %v%s\n", "\x1b[91m", err, "\x1b[0m")
		os.Exit(1)
	}

	utils.InstallClink()
	service.OptimizeServices()
	registry.OptimizeRegistry()
	utils.ClearQuickAccess()
	utils.InstallPNPM()
	utils.InstallFont()
	utils.ConfigureOhMyPosh()
	utils.InstallDraculaTheme()
	utils.InstallVSCode()

	fmt.Printf("%s[+] Xong rồi ní, code đi thôi!%s\n", "\x1b[32m", "\x1b[0m")
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
