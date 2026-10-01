package main

import (
	"syscall"
	"unsafe"
)

// enableUTF8Console switches the console codepage to UTF-8 (65001) so Chinese
// output is not garbled on legacy conhost (default 936 on zh-CN Windows).
func enableUTF8Console() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	if p := kernel32.NewProc("SetConsoleOutputCP"); p != nil {
		_, _, _ = p.Call(65001)
	}
	if p := kernel32.NewProc("SetConsoleCP"); p != nil {
		_, _, _ = p.Call(65001)
	}
}

func IsDoubleClickRun() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}
