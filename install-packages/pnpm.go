package install_packages

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type NpmPackageInfo struct {
	DistTags map[string]string `json:"dist-tags"`
}

func InstallPnpm() {
	// get latest pnpm version from npmjs
	resp, err := http.Get("https://registry.npmjs.org/@pnpm/exe")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	var pkgInfo NpmPackageInfo
	if err := json.NewDecoder(resp.Body).Decode(&pkgInfo); err != nil {
		return
	}
	version := pkgInfo.DistTags["latest"]
	if version == "" {
		return
	}

	// download pnpm
	downloadUrl := fmt.Sprintf("https://github.com/pnpm/pnpm/releases/download/v%s/pnpm-win-x64.exe", version)
	resp, err = http.Get(downloadUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	file, err := os.Create("pnpm-setup.exe")
	if err != nil {
		log.Fatal(err)
		return
	}
	io.Copy(file, resp.Body)
	file.Close()

	// install pnpm
	cmd := exec.Command("./pnpm-setup.exe", "setup")
	cmd.Run()

	// set env
	localAppData := os.Getenv("LOCALAPPDATA")
	pnpmHome := filepath.Join(localAppData, "pnpm")
	os.Setenv("PNPM_HOME", pnpmHome)

	path := os.Getenv("PATH")
	newPath := pnpmHome + ";" + path
	os.Setenv("PATH", newPath)

	// install nodejs
	cmd = exec.Command("pnpm", "env", "use", "--global", "jod")
	cmd.Run()

	// remove prefix when install
	cmd = exec.Command("pnpm", "config", "set", "save-prefix", "")
	cmd.Run()
}
