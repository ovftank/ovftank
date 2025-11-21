package utils

import (
	"archive/zip"
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
		fmt.Printf("[+] jetbrainsmono đã có sẵn, skip\n")
		return
	}

	fmt.Printf("[*] Đang cài đặt JetBrainsMono Font...\n")

	fontUrl := "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.4.0/JetBrainsMono.zip"
	fontZip := filepath.Join(os.Getenv("TEMP"), "JetBrainsMono.zip")
	fontExtractPath := filepath.Join(os.Getenv("TEMP"), "JetBrainsMono")

	fmt.Printf("[*] tải font zip từ %s\n", fontUrl)
	resp, err := http.Get(fontUrl)
	if err != nil {
		fmt.Printf("[!] tải font fail: %v\n", err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		file, err := os.Create(fontZip)
		if err != nil {
			fmt.Printf("[!] tạo file zip fail: %v\n", err)
			return
		}
		if file != nil {
			defer file.Close()
			if _, err := io.Copy(file, resp.Body); err != nil {
				fmt.Printf("[!] lưu zip fail: %v\n", err)
				return
			}
			file.Close()
			fmt.Printf("[+] tải xong: %s\n", fontZip)
		}
	}

	if err := os.MkdirAll(fontExtractPath, 0755); err != nil {
		fmt.Printf("[!] tạo folder giải nén fail: %v\n", err)
		return
	}
	fmt.Printf("[*] chuẩn bị giải nén vào %s\n", fontExtractPath)

	if err := extractFontZip(fontZip, fontExtractPath); err != nil {
		fmt.Printf("[!] giải nén zip fail: %v\n", err)
		return
	}

	filepath.Walk(fontExtractPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".ttf") || strings.HasSuffix(strings.ToLower(info.Name()), ".otf") {
			destPath := filepath.Join(fontDestination, info.Name())
			copyFile(path, destPath)
			fmt.Printf("[*] copy font %s\n", info.Name())

			fontKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.SET_VALUE)
			if err == nil {
				defer fontKey.Close()
				fontName := strings.ReplaceAll(info.Name(), ".ttf", " (TrueType)")
				fontName = strings.ReplaceAll(fontName, ".otf", " (TrueType)")
				fontKey.SetStringValue(fontName, info.Name())
				fmt.Printf("[*] update registry cho %s\n", fontName)
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
		fmt.Printf("[*] cấu hình console font\n")
	}

	os.Remove(fontZip)
	os.RemoveAll(fontExtractPath)
	fmt.Printf("[+] dọn temp font xong\n")
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

func extractFontZip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		targetPath := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		src, err := f.Open()
		if err != nil {
			return err
		}

		dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			src.Close()
			return err
		}

		if _, err := io.Copy(dst, src); err != nil {
			src.Close()
			dst.Close()
			return err
		}

		src.Close()
		dst.Close()
	}

	return nil
}
