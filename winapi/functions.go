package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                = syscall.NewLazyDLL("user32.dll")
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	getLastError          = kernel32.NewProc("GetLastError")
	getModuleHandleA      = kernel32.NewProc("GetModuleHandleA")
	registerHotkeyW       = user32.NewProc("RegisterHotKey")
	unregisterHotkey      = user32.NewProc("UnregisterHotKey")
	getMessageW           = user32.NewProc("GetMessageW")
	dispatchMessageW      = user32.NewProc("DispatchMessageW")
	messageBoxW           = user32.NewProc("MessageBoxW")
	getForegroundWindowW  = user32.NewProc("GetForegroundWindow")
	setForegroundWindow   = user32.NewProc("SetForegroundWindow")
	setWindowTextW        = user32.NewProc("SetWindowTextW")
	moveWindow            = user32.NewProc("MoveWindow")
	getWindowRect         = user32.NewProc("GetWindowRect")
	systemParametersInfoA = user32.NewProc("SystemParametersInfoA")
	getMessagePos         = user32.NewProc("GetMessagePos")
	getCursorPos          = user32.NewProc("GetCursorPos")
)

// RegisterHotkey will register the given key
//
func RegisterHotkey(id int, key rune, alt bool, ctrl bool, shift bool) bool {
	return RegisterHotVkey(id, CharToVK(key), alt, ctrl, shift)
}

// RegisterHotVkey will register a hotkey using its virtual key code
func RegisterHotVkey(id int, vkey int32, alt bool, ctrl bool, shift bool) bool {
	var modifier = 0
	if alt {
		modifier |= MOD_ALT
	}
	if shift {
		modifier |= MOD_SHIFT
	}
	if ctrl {
		modifier |= MOD_CONTROL
	}
	ret, _, _ := registerHotkeyW.Call(uintptr(0), uintptr(id), uintptr(modifier), uintptr(vkey))
	return ret != 0
}

// UnregisterHotVkey will unregister the hotkey binding with the given id
func UnregisterHotVkey(id int) bool {
	ret, _, _ := unregisterHotkey.Call(uintptr(0), uintptr(id))
	return ret != 1
}

// GetMessage retrieves a message from the windows message queue in a blocking way
func GetMessage(msg *Message) bool {
	ret, _, _ := getMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(0), uintptr(0x0), uintptr(0))
	return ret != 0
}

// DispatchMessage will forward messages to the corresponding windows
func DispatchMessage(msg *Message) bool {
	ret, _, _ := dispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

// GetMessagePos retrieves the coordinates of the last message
func GetMessagePos() (x int, y int) {
	ret, _, _ := getMessagePos.Call()
	x = int(ret & 0xFFFF)
	y = int(ret >> 16)
	return
}

func GetCursorPos() (x int, y int) {
	coords := new(Point)
	_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(coords)))
	return int(coords.x), int(coords.y)
}

// MessageBox will show a message box
func MessageBox(title string, description string) uintptr {
	ret, _, _ := messageBoxW.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(description))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(MB_YESNOCANCEL))
	return ret
}

// GetForegroundWindow retrieves the handle of the currently actrive foreground window
func GetForegroundWindow() Hwnd {
	ret, _, _ := getForegroundWindowW.Call()
	return Hwnd(ret)
}

// SetForegroundWindow brings a window to the front
func SetForegroundWindow(hwnd Hwnd) {
	setForegroundWindow.Call(uintptr(hwnd))
}

// SetWindowText sets the title bar of the window if it exists
func SetWindowText(hwnd Hwnd, title string) bool {
	ret, _, _ := setWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
	return ret != 0
}

// MoveWindow position the window with the given handle
func MoveWindow(hwnd Hwnd, x int, y int, nWidth int, nHeight int) bool {
	ret, _, _ := moveWindow.Call(uintptr(hwnd), uintptr(x), uintptr(y), uintptr(nWidth), uintptr(nHeight), uintptr(1))
	return ret != 0
}

// GetWindowRect returns then window rectangle of the given window
func GetWindowRect(hwnd Hwnd, rect *Rect) bool {
	ret, _, _ := getWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect)))
	return ret != 0
}

// GetWorkArea returns the usable size minus taskbar and or app bars for the primary screen
func GetWorkArea() (*Rect, bool) {
	var rect = new(Rect)
	ret, _, _ := systemParametersInfoA.Call(uintptr(SPI_GETWORKAREA), uintptr(0), uintptr(unsafe.Pointer(rect)), uintptr(0))
	return rect, ret != 0
}

// GetModuleHandle retrieves the handle to the given module
func GetModuleHandle(moduleName *string) uintptr {
	var name uintptr = 0
	if moduleName != nil {
		name = uintptr(unsafe.Pointer(syscall.StringBytePtr(*moduleName)))
	}
	ret, _, err := getModuleHandleA.Call(name)
	fmt.Printf("Get Module: %d, %s\n", ret, err)
	return ret
}

func ifError(condition bool, error error, prefix string) error {
	if condition {
		return fmt.Errorf(fmt.Sprintf("%s: %s\n", prefix, error.Error()))
	}
	return nil
}

// CharToVK converts a rune to a virtual key (valid for a-z)
func CharToVK(char rune) int32 {
	return char - 'a' + VK_A
}
