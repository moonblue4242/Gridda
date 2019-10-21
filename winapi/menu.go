//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package winapi

import (
	"syscall"
	"unsafe"
)

const (
	MF_POPUP     = 0x00000010
	MF_SEPARATOR = 0x00000800
	MF_CHECKED   = 0x00000008
	MF_UNCHECKED = 0x0
	WM_COMMAND   = 0x0111
)

var (
	createMenu      = user32.NewProc("CreateMenu")
	createPopupMenu = user32.NewProc("CreatePopupMenu")
	insertMenu      = user32.NewProc("InsertMenuW")
	modifyMenu      = user32.NewProc("ModifyMenuW")
	setMenu         = user32.NewProc("SetMenu")
	trackPopupMenu  = user32.NewProc("TrackPopupMenu")
	destroyMenu     = user32.NewProc("DestroyMenu")
)

type (
	// HMenu describes a handle to a menu
	HMenu uintptr
)

// Menu allows controlling a windows menu
type Menu interface {
	Set(hwnd Hwnd) error
	Insert(position uint, flags uint, id uint, label string) error
	Modify(position uint, flags uint, id uint, label string) error
	InsertSeparator() error
	InsertAtParent(position uint, flags uint, parent Menu, label string) error
	Destroy() error
	handle() HMenu
}

// PopupMenu is a menu which can be opened anywhere on the screen like e.g. a context menu
type PopupMenu interface {
	Menu
	Track(flags uint, x int, y int, hwnd Hwnd) error
}

type menu struct {
	handleMenu HMenu
}

// NewMenu creates a new Windows menu
func NewMenu() (Menu, error) {
	ret, _, error := createMenu.Call()
	return &menu{handleMenu: HMenu(ret)}, ifError(ret == 0, error, "CreateMenu")
}

// NewPopupMenu creates a new Windows popup menu
func NewPopupMenu() (PopupMenu, error) {
	ret, _, error := createPopupMenu.Call()
	return &menu{handleMenu: HMenu(ret)}, ifError(ret == 0, error, "CreatePopupMenu")
}

// SetMenu adds a menu to a window
func (menu *menu) Set(hwnd Hwnd) error {
	_, _, err := setMenu.Call(uintptr(hwnd), uintptr(menu.handleMenu))
	return err //ifError(ret == 0, error, "SetMenu")
}

// InsertAtParent inserts a new menu item to the top
func (menu *menu) Insert(position uint, flags uint, id uint, label string) error {
	ret, _, err := insertMenu.Call(uintptr(menu.handleMenu), uintptr(position), uintptr(flags), uintptr(id), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(label))))
	return ifError(ret != 1, err, "InsertMenu")
}

// InsertSeparator will insert a horizontal separator
func (menu *menu) InsertSeparator() error {
	ret, _, err := insertMenu.Call(uintptr(menu.handleMenu), uintptr(0), uintptr(MF_SEPARATOR), uintptr(1), uintptr(0))
	return ifError(ret != 1, err, "InsertSeparator")
}

// InsertAtParent inserts a new menu item to the top
func (menu *menu) InsertAtParent(position uint, flags uint, parent Menu, label string) error {
	ret, _, err := insertMenu.Call(uintptr(menu.handleMenu), uintptr(position), uintptr(flags), uintptr(parent.handle()), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(label))))
	return ifError(ret != 1, err, "InsertAtParent")
}

// Modify updates a menu item
func (menu *menu) Modify(position uint, flags uint, id uint, label string) error {
	ret, _, err := modifyMenu.Call(uintptr(menu.handleMenu), uintptr(position), uintptr(flags), uintptr(id), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(label))))
	return ifError(ret != 1, err, "Modify")
}

// Track will show a popup menu
func (menu *menu) Track(flags uint, x int, y int, hwnd Hwnd) error {
	ret, _, err := trackPopupMenu.Call(uintptr(menu.handleMenu), uintptr(flags), uintptr(x), uintptr(y), uintptr(0), uintptr(hwnd), uintptr(0))
	return ifError(ret != 1, err, "TrackMenu")
}

// Destroy will remove the menu item
func (menu *menu) Destroy() error {
	ret, _, err := destroyMenu.Call(uintptr(menu.handleMenu))
	return ifError(ret != 1, err, "DestroyMenu")
}

func (menu *menu) handle() HMenu {
	return menu.handleMenu
}
