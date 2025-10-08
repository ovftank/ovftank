package utils

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func QuickAccess() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	shell, err := oleutil.CreateObject("Shell.Application")
	if err != nil {
		return
	}
	defer shell.Release()

	shellDispatch, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return
	}
	defer shellDispatch.Release()

	quickAccess, err := oleutil.CallMethod(shellDispatch, "Namespace", "shell:::{679F85CB-0220-4080-B29B-5540CC05AAB6}")
	if err != nil {
		return
	}
	quickAccessObj := quickAccess.ToIDispatch()
	defer quickAccessObj.Release()

	items, err := oleutil.CallMethod(quickAccessObj, "Items")
	if err != nil {
		return
	}
	itemsObj := items.ToIDispatch()
	defer itemsObj.Release()

	count, err := oleutil.GetProperty(itemsObj, "Count")
	if err != nil {
		return
	}

	for i := 0; i < int(count.Val); i++ {
		item, err := oleutil.CallMethod(itemsObj, "Item", i)
		if err != nil {
			continue
		}
		itemObj := item.ToIDispatch()

		verbs, err := oleutil.CallMethod(itemObj, "Verbs")
		if err != nil {
			itemObj.Release()
			continue
		}
		verbsObj := verbs.ToIDispatch()
		defer verbsObj.Release()

		verbCount, err := oleutil.GetProperty(verbsObj, "Count")
		if err != nil {
			itemObj.Release()
			continue
		}

		for j := 0; j < int(verbCount.Val); j++ {
			verb, err := oleutil.CallMethod(verbsObj, "Item", j)
			if err != nil {
				continue
			}
			verbObj := verb.ToIDispatch()

			name, err := oleutil.GetProperty(verbObj, "Name")
			if err == nil && name.ToString() == "Unpin from Quick access" {
				oleutil.CallMethod(verbObj, "DoIt")
			}
			verbObj.Release()
		}
		itemObj.Release()
	}
}

func ClearQuickAccess() {
	fmt.Printf("[*] Đang xóa Quick Access...\n")
	QuickAccess()
	fmt.Printf("[+] Đã xóa Quick Access\n")
}
