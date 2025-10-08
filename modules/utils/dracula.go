package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallDraculaTheme() {
	// Check if theme already installed
	draculaThemePath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows Terminal", "themes", "dracula.json")
	if _, err := os.Stat(draculaThemePath); err == nil {
		fmt.Printf("[+] Dracula theme đã cài đặt... bỏ qua\n")
		return
	}

	fmt.Printf("[*] Đang cài đặt Dracula theme...\n")
	colorToolUrl := "https://raw.githubusercontent.com/waf/dracula-cmd/master/dist/ColorTool.zip"
	colorToolZip := filepath.Join(os.Getenv("TEMP"), "ColorTool.zip")
	colorToolPath := filepath.Join(os.Getenv("TEMP"), "ColorTool")

	if _, err := os.Stat(colorToolZip); err == nil {
		os.Remove(colorToolZip)
	}
	if _, err := os.Stat(colorToolPath); err == nil {
		os.RemoveAll(colorToolPath)
	}

	downloadFile(colorToolUrl, colorToolZip)

	extractZip(colorToolZip, colorToolPath)

	actualColorToolPath := filepath.Join(colorToolPath, "ColorTool")
	var colorToolExe, installFolder string

	if _, err := os.Stat(actualColorToolPath); err == nil {
		colorToolExe = filepath.Join(actualColorToolPath, "ColorTool.exe")
		installFolder = filepath.Join(actualColorToolPath, "install")
	} else {
		colorToolExe = filepath.Join(colorToolPath, "ColorTool.exe")
		installFolder = filepath.Join(colorToolPath, "install")
	}

	if _, err := os.Stat(colorToolExe); err == nil {
		shortcutPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Windows PowerShell")

		if _, err := os.Stat(filepath.Join(installFolder, "Windows PowerShell.lnk")); err == nil {
			copyFile(filepath.Join(installFolder, "Windows PowerShell.lnk"), filepath.Join(shortcutPath, "Windows PowerShell.lnk"))
		}
		if _, err := os.Stat(filepath.Join(installFolder, "Windows PowerShell (x86).lnk")); err == nil {
			copyFile(filepath.Join(installFolder, "Windows PowerShell (x86).lnk"), filepath.Join(shortcutPath, "Windows PowerShell (x86).lnk"))
		}

		regFile := filepath.Join(installFolder, "Remove Default Console Overrides.reg")
		if _, err := os.Stat(regFile); err == nil {
			exec.Command("reg.exe", "import", regFile).Run()
		}

		themeFile := filepath.Join(installFolder, "Dracula-ColorTool.itermcolors")
		if _, err := os.Stat(themeFile); err == nil {
			exec.Command(colorToolExe, "-b", themeFile).Run()
		}
	}

	os.Remove(colorToolZip)
	os.RemoveAll(colorToolPath)
}

func extractZip(src, dest string) {
	r, _ := zip.OpenReader(src)
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		outFile, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		rc, _ := f.Open()
		io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}
}
