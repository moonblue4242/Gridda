package ui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/moonblue4242/Gridda/cmds"

	"github.com/moonblue4242/Gridda/winapi"
)

const (
	iconCount        = 6
	trayMsgID        = 1111
	trayExitCmd      = 4711
	traySeparator    = trayExitCmd + 1
	trayGridOneCmd   = traySeparator + 1
	trayGridTwoCmd   = trayGridOneCmd + 1
	trayGridThreeCmd = trayGridOneCmd + 2
	trayGridFourCmd  = trayGridOneCmd + 3
	trayGridFiveCmd  = trayGridOneCmd + 4
	trayGridSixCmd   = trayGridOneCmd + 5
)

// TrayIcon is an abstraction for the tray icon of the application
type TrayIcon interface {
	OnCommand(msg *winapi.Message)
	SetGrid(config *cmds.Config, index int)
	Shutdown()
}

// TrayIconActions defines all the actions which can be triggered by the user via the tray icon
type TrayIconActions interface {
	OnShowApplication()
	OnExitApplication()
	OnGridSelect(index int)
}

type trayIcon struct {
	hwnd         winapi.Hwnd
	notification winapi.Notification
	icons        [iconCount]winapi.Icon
	actions      TrayIconActions
	menu         winapi.PopupMenu
	selected     int
}

// NewTrayIcon creates a new tray for Grida encapsulating all tray related features
func NewTrayIcon(hwnd winapi.Hwnd, actions TrayIconActions, config *cmds.Config) TrayIcon {
	tray := &trayIcon{hwnd: hwnd, actions: actions}
	tray.setupTrayIcon()
	tray.setupMenu(config, 0)
	return tray
}

// OnCommand must be called by the event loop whenever a WM_COMMAND is received
func (trayIcon *trayIcon) OnCommand(msg *winapi.Message) {
	switch msg.WParam {
	case trayExitCmd:
		trayIcon.actions.OnExitApplication()
	case trayGridOneCmd, trayGridTwoCmd, trayGridThreeCmd, trayGridFourCmd, trayGridFiveCmd, trayGridSixCmd:
		trayIcon.actions.OnGridSelect(int(msg.WParam) - trayGridOneCmd)
	}
	return
}

// SetGrid selects the active grid to show as icon and as menu entry via checked item
func (trayIcon *trayIcon) SetGrid(config *cmds.Config, index int) {
	trayIcon.setupMenu(config, index)
	trayIcon.notification.SetIcon(trayIcon.icons[index])
}

func (trayIcon *trayIcon) setupTrayIcon() error {
	var err error
	// get icons
	for i := 1; i <= iconCount; i++ {
		trayIcon.icons[i-1], err = winapi.NewIcon(strconv.Itoa(i) + ".ico")
		if err != nil {
			log.Fatalf("Unable to load icon #%d: %s\n", i, err)
		}
	}
	// setup notification icon and basic UI
	trayIcon.notification, err = winapi.AddNotification(trayIcon.hwnd, 0, trayMsgID, trayIcon.icons[0])
	fmt.Printf("Notification Icon: %s\n", err)
	// show window on double click
	winapi.AddMessageHandler(trayIcon.hwnd, winapi.NewMessageHandler(trayMsgID, func(hwnd winapi.Hwnd, wParam winapi.Wparam, lParam winapi.LParam) bool {
		if lParam == winapi.WM_LBUTTONDBLCLK {
			x, y := winapi.GetMessagePos()
			fmt.Printf("Notification double click received (x:%d, y:%d)\n", x, y)
			trayIcon.actions.OnShowApplication()
		} else if lParam == winapi.WM_RBUTTONUP {
			x, y := winapi.GetCursorPos()
			winapi.SetForegroundWindow(trayIcon.hwnd)
			trayIcon.menu.Track(0, x, y, trayIcon.hwnd)
			fmt.Printf("Notification right up received (x:%d, y:%d)\n", x, y)
		}
		return true
	}))
	return err
}

func (trayIcon *trayIcon) setupMenu(config *cmds.Config, selected int) {
	var err error
	if trayIcon.menu != nil {
		trayIcon.menu.Destroy()
	}
	trayIcon.menu, err = winapi.NewPopupMenu()
	if err != nil {
		log.Fatalf("Unable to create popup menu: %s\n", err)
	}
	for i := 0; i < iconCount && i < len(config.Grids); i++ {
		var flags uint = 0
		if i == selected {
			flags |= winapi.MF_CHECKED
		}
		trayIcon.menu.Insert(0, flags, trayGridOneCmd+uint(i), strconv.Itoa(i+1)+" - "+config.Grids[i].Name)
	}
	trayIcon.menu.InsertSeparator()
	trayIcon.menu.Insert(0, 0, trayExitCmd, "Exit")
}

// Shutdown will destroy and free all tray related resources
func (trayIcon *trayIcon) Shutdown() {
	if err := trayIcon.notification.Remove(); err != nil {
		log.Println(err)
	}
}
