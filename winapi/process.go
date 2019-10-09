package winapi

import (
	"fmt"
	"unsafe"
)

const (
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_READ           = 0x0010
	READ_CONTROL              = 0x00020000
)

var (
	getWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	getModuleBaseNameA       = kernel32.NewProc("K32GetModuleBaseNameA")
	openProcess              = kernel32.NewProc("OpenProcess")
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
