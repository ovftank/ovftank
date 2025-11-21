package service

import (
	"fmt"

	"golang.org/x/sys/windows/svc/mgr"
)

func DisableSearch() {
	fmt.Printf("[*] Đang tắt Windows Search...\n")
	disableService("WSearch")
	fmt.Printf("[+] Đã tắt Windows Search\n")
}

func DisableTime() {
	fmt.Printf("[*] Đang tắt Windows Time...\n")
	disableService("w32time")
	fmt.Printf("[+] Đã tắt Windows Time\n")
}

func DisableUpdate() {
	fmt.Printf("[*] Đang tắt Windows Update...\n")
	disableService("wuauserv")
	fmt.Printf("[+] Đã tắt Windows Update\n")
}

func DisablePrint() {
	fmt.Printf("[*] Đang tắt Print Spooler...\n")
	disableService("Spooler")
	fmt.Printf("[+] Đã tắt Print Spooler\n")
}

func DisablePlugPlay() {
	fmt.Printf("[*] Đang tắt Plug and Play...\n")
	disableService("PlugPlay")
	fmt.Printf("[+] Đã tắt Plug and Play\n")
}

func disableService(serviceName string) {
	m, err := mgr.Connect()
	if err != nil {
		return
	}
	defer m.Disconnect()

	service, err := m.OpenService(serviceName)
	if err != nil {
		return
	}
	defer service.Close()

	config, err := service.Config()
	if err != nil {
		return
	}
	config.StartType = mgr.StartDisabled
	service.UpdateConfig(config)
}

func OptimizeServices() {
	fmt.Printf("[*] Đang tắt các services không cần thiết...\n")
	DisableSearch()
	DisableTime()
	DisableUpdate()
	DisablePrint()
	DisablePlugPlay()
}
