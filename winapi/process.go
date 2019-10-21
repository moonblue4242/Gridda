//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package winapi

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_READ           = 0x0010
	READ_CONTROL              = 0x00020000
)

var (
	hooky                    = syscall.NewLazyDLL("hooky.dll")
	getWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	getModuleBaseNameA       = kernel32.NewProc("K32GetModuleBaseNameA")
	openProcess              = kernel32.NewProc("OpenProcess")
	loadLibrary              = kernel32.NewProc("LoadLibraryW")
	getProcAddr              = kernel32.NewProc("GetProcAddress")
	setWindowsHookEx         = user32.NewProc("SetWindowsHookExW")
	unhookWindowsHookEx      = user32.NewProc("UnhookWindowsHookEx")
	registerWindowMessage    = user32.NewProc("RegisterWindowMessageA")
)

// GetWindowThreadProcessID returns the process id of the window owning process
func GetWindowThreadProcessID(hwnd Hwnd) int32 {
	var processID int32
	getWindowThreadProcessID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processID)))
	return processID
}

// GetModuleBaseName will return the base name (e.g. the name of the executable) of the given process
func GetModuleBaseName(processHandle uintptr) string {
	var data [512]byte
	var result string = ""
	ret, _, err := getModuleBaseNameA.Call(processHandle, uintptr(0), uintptr(unsafe.Pointer(&data)), uintptr(len(data)))
	fmt.Printf("GetWindowModuleFileName: %d, %s\n", ret, err)
	if ret != 0 {
		result = fmt.Sprintf("%s", data[:ret])
	}
	return result
}

// OpenProcess retrieves a handle to a process in read only mod with PROCESS_VM_READ, PROCESS_QUERY_INFORMATION and READ_CONTROL set
func OpenProcess(processID int32) (uintptr, error) {
	ret, _, err := openProcess.Call(uintptr(READ_CONTROL|PROCESS_VM_READ|PROCESS_QUERY_INFORMATION), uintptr(0), uintptr(processID))
	fmt.Printf("Openprocess: %s\n", err)
	return ret, ifError(ret == 0, err, "OpenProcess")
}

// AddHook will register a hook
func AddHook() uintptr {
	var hook uintptr

	hModule, _, err := loadLibrary.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("hooky.dll"))))
	// hModule := hooky.Handle()
	if hModule != 0 {
		ret, _, err := getProcAddr.Call(uintptr(hModule), uintptr(unsafe.Pointer(syscall.StringBytePtr("fnCbtProc"))))
		fmt.Printf("GetProcAddr: %d, %s\n", ret, err)
		hook, _, err = setWindowsHookEx.Call(uintptr(5), uintptr(ret), uintptr(hModule), uintptr(0))
		fmt.Printf("SetWindowsHookEx: %d, %s\n", hook, err)
	} else {
		log.Panicln(err)
	}
	return hook
}

func RemoveHook(hook uintptr) {
	ret, _, err := unhookWindowsHookEx.Call(hook)
	fmt.Printf("Unhook: %d, %s\n", ret, err)
}

// RegisterWindowMessage will register a global message identifier under the given id
func RegisterWindowMessage(id string) uint {
	ret, _, err := registerWindowMessage.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr(id))))
	if ret == 0 {
		log.Panicf("RegisterWindowMessage: %s\n", err)
	}
	return uint(ret)

}
