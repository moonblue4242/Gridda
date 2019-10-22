package winapi

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	HCBT_ACTIVATE     = 5
	HCBT_CLICKSKIPPED = 6
	HCBT_CREATEWND    = 3
	HCBT_DESTROYWND   = 4
	HCBT_KEYSKIPPED   = 7
	HCBT_MINMAX       = 1
	HCBT_MOVESIZE     = 0
	HCBT_SETFOCUS     = 9

	WH_CBT = 5
)

var (
	hooky                 = syscall.NewLazyDLL("hooky.dll")
	hookyCbtProc          = hooky.NewProc("fnCbtProc")
	setWindowsHookEx      = user32.NewProc("SetWindowsHookExW")
	unhookWindowsHookEx   = user32.NewProc("UnhookWindowsHookEx")
	registerWindowMessage = user32.NewProc("RegisterWindowMessageA")
)

// Hook defines a windows system hook
type Hook uintptr

// AddCbtHook will register a hook
func AddCbtHook() Hook {
	hModule := hooky.Handle()
	hook, _, err := setWindowsHookEx.Call(uintptr(WH_CBT), hookyCbtProc.Addr(), uintptr(hModule), uintptr(0))
	if hook == 0 {
		log.Panicf("AddHook: %s\n", err)
	}
	return Hook(hook)
}

// RemoveHook will remove the given hook
func RemoveHook(hook Hook) {
	ret, _, err := unhookWindowsHookEx.Call(uintptr(hook))
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
