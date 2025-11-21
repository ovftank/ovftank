package registry

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func DisableSearch() {
	const parentKeyPath = `SOFTWARE\Policies\Microsoft\Windows`
	const searchSubKey = "Windows Search"
	const valueName = "DisableSearch"

	winKey, err := registry.OpenKey(registry.LOCAL_MACHINE, parentKeyPath, registry.CREATE_SUB_KEY)
	if err != nil {
		return
	}
	defer winKey.Close()

	searchKey, _, err := registry.CreateKey(winKey, searchSubKey, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer searchKey.Close()
	searchKey.SetDWordValue(valueName, 1)
}

func HideSearch() {
	const keyPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\Search`
	const valueName = "SearchboxTaskbarMode"
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer key.Close()
	key.SetDWordValue(valueName, 0)
}

func EnableReg() {
	hklmKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.SET_VALUE)
	if err != nil {
		hklmKey, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.SET_VALUE)
		if err != nil {
			return
		}
	}
	defer hklmKey.Close()
	hklmKey.SetDWordValue("DisableRegistryTools", 0)

	hkcuKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.SET_VALUE)
	if err != nil {
		hkcuKey, _, err = registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.SET_VALUE)
		if err != nil {
			return
		}
	}
	defer hkcuKey.Close()
	hkcuKey.SetDWordValue("DisableRegistryTools", 0)
}

func OptimizeExplorer() {
	explorerAdvancedKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, registry.SET_VALUE)
	if err == nil {
		defer explorerAdvancedKey.Close()
		explorerAdvancedKey.SetDWordValue("LaunchTo", 1)
		explorerAdvancedKey.SetDWordValue("Hidden", 1)
		explorerAdvancedKey.SetDWordValue("HideFileExt", 0)
	}

	explorerKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer`, registry.SET_VALUE)
	if err == nil {
		defer explorerKey.Close()
		explorerKey.SetDWordValue("ShowRecent", 0)
		explorerKey.SetDWordValue("ShowFrequent", 0)
	}
}

func OptimizeAccess() {
	stickyKeysKey, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Accessibility\StickyKeys`, registry.SET_VALUE)
	if err == nil {
		defer stickyKeysKey.Close()
		stickyKeysKey.SetStringValue("Flags", "506")
	}

	toggleKeysKey, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Accessibility\ToggleKeys`, registry.SET_VALUE)
	if err == nil {
		defer toggleKeysKey.Close()
		toggleKeysKey.SetStringValue("Flags", "58")
	}

	keyboardResponseKey, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Accessibility\Keyboard Response`, registry.SET_VALUE)
	if err == nil {
		defer keyboardResponseKey.Close()
		keyboardResponseKey.SetStringValue("Flags", "122")
	}
}

func OptimizeVisual() {
	visualEffectsKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, registry.SET_VALUE)
	if err == nil {
		defer visualEffectsKey.Close()
		visualEffectsKey.SetDWordValue("VisualFXSetting", 2)
	}
}

func OptimizeDns() {
	dnsKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Dnscache\Parameters`, registry.SET_VALUE)
	if err == nil {
		defer dnsKey.Close()
		dnsKey.SetDWordValue("MaxCacheTtl", 86400)
		dnsKey.SetDWordValue("MaxNegativeCacheTtl", 0)
	}
}
func DarkTheme() {
	personalizeKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.SET_VALUE)
	if err == nil {
		defer personalizeKey.Close()
		personalizeKey.SetDWordValue("EnableTransparency", 0)
		personalizeKey.SetDWordValue("AppsUseLightTheme", 0)
		personalizeKey.SetDWordValue("SystemUsesLightTheme", 0)
	}
}

func ThisPC() {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{088e3905-0323-4b02-9826-5d99428e115f}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{24ad3ad4-a569-4530-98e1-ab02f9417aa8}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{3dfdf296-dbec-4fb4-81d1-6a3438bcf4de}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{f86fa3ab-70d2-4fc7-9c99-fcbf05467f3a}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{d3162b92-9365-467a-956b-92703aca08af}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{B4BFCC3A-DB2C-424C-B029-7FE99A87C641}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{088e3905-0323-4b02-9826-5d99428e115f}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{24ad3ad4-a569-4530-98e1-ab02f9417aa8}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{3dfdf296-dbec-4fb4-81d1-6a3438bcf4de}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{f86fa3ab-70d2-4fc7-9c99-fcbf05467f3a}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{d3162b92-9365-467a-956b-92703aca08af}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}

	key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\{B4BFCC3A-DB2C-424C-B029-7FE99A87C641}`, registry.SET_VALUE)
	if err == nil {
		key.Close()
	}
}

func OptimizeRegistry() {
	fmt.Printf("[*] Đang cấu hình registry...\n")
	EnableReg()
	DisableSearch()
	HideSearch()
	OptimizeExplorer()
	OptimizeAccess()
	OptimizeVisual()
	OptimizeDns()
	DarkTheme()
	ThisPC()
	fmt.Printf("[+] Đã cấu hình registry xong\n")
}
