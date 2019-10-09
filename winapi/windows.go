package winapi

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	defaultWindowClass = "GRIDDA_DEFAULT_WINDOW_CLASS"

	wS_BORDER        = 0x00800000
	wS_CAPTION       = 0x00C00000
	wS_SYSMENU       = 0x00080000
	wS_VISIBLE       = 0x10000000
	wS_EX_NOACTIVATE = 0x08000000
)

var (
	handlers        messageHandlers = make(messageHandlers)
	registerClassW                  = user32.NewProc("RegisterClassW")
	createWindowExW                 = user32.NewProc("CreateWindowExW")
	defWindowProcW                  = user32.NewProc("DefWindowProcW")
	showWindow                      = user32.NewProc("ShowWindow")
	getClientRect                   = user32.NewProc("GetClientRect")
	enumWindows                     = user32.NewProc("EnumWindows")
	isWindow                        = user32.NewProc("IsWindow")
	getWindowText                   = user32.NewProc("GetWindowTextA")
)

// MessageHandler can be used to add a custom handler for a specific window or application message to
// the window event loop
type MessageHandler struct {
	onEvent func(hwnd Hwnd, wParam Wparam, lParam LParam) bool
	msgID   uint
}

// NewMessageHandler creates a new handler object initializing it for correct use
func NewMessageHandler(msgID uint, onEvent func(hwnd Hwnd, wParam Wparam, lParam LParam) bool) *MessageHandler {
	return &MessageHandler{msgID: msgID, onEvent: onEvent}
}

// AddMessageHandler will add a message handler to the given window
func AddMessageHandler(hwnd Hwnd, messageHandler *MessageHandler) {
	handlers.add(hwnd, messageHandler)
}

// type mm map[uint]MessageHandler
type messageHandlers map[Hwnd]map[uint]*MessageHandler

func (messageHandlers messageHandlers) add(hwnd Hwnd, messageHandler *MessageHandler) {
	item, ok := messageHandlers[hwnd]
	if !ok {
		item = make(map[uint]*MessageHandler)
		messageHandlers[hwnd] = item
	}
	item[messageHandler.msgID] = messageHandler
}

type winProc func(hwnd Hwnd, msg uint, wParam Wparam, lParam LParam) int

// WinEnumProc defines a callback function called during window enumeration
type WinEnumProc func(hwnd Hwnd, lParam LParam) uintptr

type wNDCLASSEX struct {
	Sizeof        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	iconSm        uintptr
}

func windowCallback(hwnd Hwnd, msg uint, wParam Wparam, lParam LParam) int {
	var done = false

	handler, ok := handlers[hwnd][msg]
	if ok {
		done = done || handler.onEvent(hwnd, wParam, lParam)
	}

	var ret uintptr = uintptr(1)
	if !done {
		ret, _, _ = defWindowProcW.Call(uintptr(hwnd), uintptr(msg), uintptr(wParam), uintptr(lParam))
	}
	return int(ret)
}

// registerWindowClass will create a new default windows class
func registerWindowClass() error {
	windowClass := new(wNDCLASSEX)
	windowClass.style = 0x0088
	windowClass.lpfnWndProc = syscall.NewCallback(windowCallback)
	windowClass.cbClsExtra = 0
	windowClass.cbWndExtra = 0
	windowClass.hInstance = GetModuleHandle(nil)
	windowClass.lpszClassName = (*uint16)(unsafe.Pointer(syscall.StringToUTF16Ptr(defaultWindowClass)))
	rcRet, _, err := registerClassW.Call(uintptr(unsafe.Pointer(windowClass)))
	fmt.Printf("registerWindowClass: %d\n", rcRet)
	return ifError(rcRet == 0, err, "registerWindowClass")
}

// CreateInactiveWindow will create a basic window which is not active by default and not visible inside the taskbar
// it will however still receive events, like e.g. shell_notifiy icon event
func CreateInactiveWindow(name string, x int, y int, width int, height int, messageHandlers ...*MessageHandler) (Hwnd, error) {
	// create window
	hwnd, _, err := createWindowExW.Call(
		uintptr(0), //dwExStyle
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(defaultWindowClass))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))),
		uintptr(0x00C00000|0x00020000|0x00080000), // dwStyle
		uintptr(x),      // x
		uintptr(y),      // y
		uintptr(width),  // widht
		uintptr(height), // height
		uintptr(0),
		uintptr(0),
		GetModuleHandle(nil),
		uintptr(0),
	)
	// add all handles
	for _, handler := range messageHandlers {
		AddMessageHandler(Hwnd(hwnd), handler)
	}
	fmt.Printf("CreateInactiveWindow: %d\n", hwnd)
	return Hwnd(hwnd), ifError(hwnd == 0, err, "CreateInactiveWindow")
}

// ShowWindow will active the window specified by the window handle by bringing it to the foreground
func ShowWindow(hwnd Hwnd) error {
	ret, _, err := showWindow.Call(uintptr(hwnd), 5)
	fmt.Printf("ShowWindow: %d, %s\n", ret, err)
	return ifError(ret == uintptr(hwnd), err, "ShowWindow")
}

// GetClientRect returns the client area of a window
func GetClientRect(hwnd Hwnd) (*Rect, error) {
	var rect Rect
	ret, _, err := getClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	fmt.Printf("GetClientRect: %d, %s\n", ret, err)
	return &rect, ifError(ret == uintptr(hwnd), err, "GetClientRect")
}

var modules = make(map[string]string)

// EnumWindows enumerates all active windows
func EnumWindows(callback WinEnumProc) error {
	ret, _, err := enumWindows.Call(uintptr(syscall.NewCallback(callback)), uintptr(4711))
	fmt.Printf("EnumWIndows: %d %s\n", ret, err)
	return ifError(ret != 1, err, "EnumWindows")
}

// GetWindowText will return the title of the window
func GetWindowText(hwnd Hwnd) (string, error) {
	var data [256]byte
	var result string = ""
	ret, _, err := getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&data)), uintptr(len(data)))
	if ret != 0 {
		result = fmt.Sprintf("%s", data[:ret])
	}
	return result, ifError(ret != 0, err, "GetWindowText")
}

func init() {
	if err := registerWindowClass(); err != nil {
		log.Fatalln(err)
	}
}
