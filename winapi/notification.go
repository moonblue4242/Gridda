//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	NIM_ADD    = 0
	NIM_MODIFY = 1
	NIM_DELETE = 2
)

var (
	shell32           = syscall.NewLazyDLL("Shell32.dll")
	shell_NotifyIconA = shell32.NewProc("Shell_NotifyIconW")
	loadImageA        = user32.NewProc("LoadImageA")
)

type gUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type notificationData struct {
	size            int32
	hwnd            Hwnd
	uID             uint32
	uFlags          uint32
	callbackMessage uint32
	hicon           uintptr
	szTip           [64]uint16
	dwState         uint32
	dwStateMask     uint32
	szInfo          [256]uint16
	uVersion        int32
	szInfoTitle     [64]uint16
	dwInfoFlags     uint32
	guidItem        gUID
}

// Notification provides functionality to control a motification icon
type Notification interface {
	SetIcon(icon Icon) error
	Remove() error
}

type notification struct {
	hwnd        Hwnd
	id          uint32
	callbackMsg int
}

// AddNotification adds a notification element to the notification area
func AddNotification(hwnd Hwnd, id int, callbackMsg int, icon Icon) (Notification, error) {

	notification := &notification{hwnd: hwnd, id: uint32(id), callbackMsg: callbackMsg}

	nd := notification.newNotificationData(icon)
	ncRet, _, err := shell_NotifyIconA.Call(uintptr(NIM_ADD), uintptr(unsafe.Pointer(nd)))

	fmt.Printf("Shell NotifyIcon: %d, %s\n", ncRet, err)

	return notification, ifError(ncRet == 1, err, "AddNotification")
}

// SetIcon updates the icon of the notification item
func (notification *notification) SetIcon(icon Icon) error {
	nd := notification.newNotificationData(icon)
	ncRet, _, err := shell_NotifyIconA.Call(uintptr(NIM_MODIFY), uintptr(unsafe.Pointer(nd)))
	return ifError(ncRet == 1, err, "SetIcon")
}

// RemoveNotification will remove the notification icon defined by the window and its window specific id
func (notification *notification) Remove() error {
	nd := new(notificationData)
	nd.size = int32(unsafe.Sizeof(nd))
	nd.hwnd = Hwnd(notification.hwnd)
	nd.uID = notification.id
	ncRet, _, err := shell_NotifyIconA.Call(uintptr(NIM_DELETE), uintptr(unsafe.Pointer(nd)))
	return ifError(ncRet == 1, err, "Remove Notification")
}

func (notification *notification) newNotificationData(icon Icon) *notificationData {
	nd := new(notificationData)
	nd.size = int32(unsafe.Sizeof(nd))
	nd.uFlags = 0x00000003
	nd.hicon = icon.handle()
	nd.hwnd = Hwnd(notification.hwnd)

	nd.callbackMessage = uint32(notification.callbackMsg)
	return nd
}
