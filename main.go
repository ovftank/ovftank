package main

import (
	"fmt"

	"windows-dev-setup/modules/utils"
)

func main() {

	fmt.Printf("\n====================================\n")
	fmt.Printf(">>> WINDOWS DEV SETUP <<<\n")
	fmt.Printf("====================================\n\n")

	// var gitName, gitEmail string
	// fmt.Printf("[INPUT] Nhập username Git: ")
	// fmt.Scanln(&gitName)
	// fmt.Printf("[INPUT] Nhập email Git: ")
	// fmt.Scanln(&gitEmail)

	// if err := chocolatey.InstallDevSetup(gitName, gitEmail); err != nil {
	// 	fmt.Printf("[!] Lỗi cài đặt: %v\n", err)
	// 	os.Exit(1)
	// }

	// utils.InstallClink()
	// service.OptimizeServices()
	// registry.OptimizeRegistry()
	// utils.ClearQuickAccess()
	// if _, err := utils.InstallPNPM(); err != nil {
	// 	fmt.Printf("[!] pnpm lỗi: %v\n", err)
	// }
	// utils.InstallFont()
	utils.ConfigureOhMyPosh()
	// utils.InstallDraculaTheme()
	// utils.InstallVSCode()
	// utils.InstallOpenKey()

	fmt.Printf("[+] Xong rồi ní, code đi thôi!\n")
	fmt.Printf("[*] nhấn Enter để thoát...")
	var exitInput string
	fmt.Scanln(&exitInput)
}
