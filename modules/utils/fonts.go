package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func InstallFont() {
	fontDestination := filepath.Join(os.Getenv("windir"), "Fonts")
	sampleFont := filepath.Join(fontDestination, "JetBrainsMonoNerdFont-Regular.ttf")

	if _, err := os.Stat(sampleFont); err == nil {
		fmt.Printf("[+] JetBrainsMono Font đã cài đặt... bỏ qua\n")
		return
	}

	fmt.Printf("[*] Đang cài đặt JetBrainsMono Font...\n")

	fontUrl := "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.4.0/JetBrainsMono.zip"
	fontZip := filepath.Join(os.Getenv("TEMP"), "JetBrainsMono.zip")
	fontExtractPath := filepath.Join(os.Getenv("TEMP"), "JetBrainsMono")

	resp, _ := http.Get(fontUrl)
	if resp != nil {
		defer resp.Body.Close()
		file, _ := os.Create(fontZip)
		if file != nil {
			defer file.Close()
			io.Copy(file, resp.Body)
			file.Close()
		}
	}

	os.MkdirAll(fontExtractPath, 0755)

	filepath.Walk(fontExtractPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".ttf") || strings.HasSuffix(strings.ToLower(info.Name()), ".otf") {
			destPath := filepath.Join(fontDestination, info.Name())
			copyFile(path, destPath)

			fontKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.SET_VALUE)
			if err == nil {
				defer fontKey.Close()
				fontName := strings.ReplaceAll(info.Name(), ".ttf", " (TrueType)")
				fontName = strings.ReplaceAll(fontName, ".otf", " (TrueType)")
				fontKey.SetStringValue(fontName, info.Name())
			}
		}
		return nil
	})

	consoleKey, err := registry.OpenKey(registry.CURRENT_USER, `Console`, registry.SET_VALUE)
	if err == nil {
		defer consoleKey.Close()
		consoleKey.SetStringValue("FaceName", "JetBrainsMono NFM")
		consoleKey.SetDWordValue("FontSize", 0x140000)
		consoleKey.SetDWordValue("WindowAlpha", 242)
	}

	os.Remove(fontZip)
	os.RemoveAll(fontExtractPath)
}

func copyFile(src, dst string) {
	source, _ := os.Open(src)
	if source == nil {
		return
	}
	defer source.Close()

	destination, _ := os.Create(dst)
	if destination == nil {
		return
	}
	defer destination.Close()

	io.Copy(destination, source)
}
