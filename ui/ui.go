//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package ui

import (
	"fmt"

	"github.com/moonblue4242/Gridda/cmds"
	"github.com/moonblue4242/Gridda/winapi"
)

const (
	hookyMsgID = "HOOKY"
)

// UI is the base instance for all ui related tasks
type UI struct {
	MainWindow winapi.Hwnd
	tray       TrayIcon
	looping    bool
	onHotkey   func(msg *winapi.Message)
	actions    Actions
	hook       winapi.Hook
}

// Actions describe all actions which can be triggered by the UI
type Actions interface {
	TrayIconActions
	HookActions
}

// HookActions describe the actions which can be triggered by the registered hook
type HookActions interface {
	OnActivate(hwnd winapi.Hwnd)
	OnFocus(hwnd winapi.Hwnd)
}

// New creates a new ui instance
func New(actions Actions, config *cmds.Config, onHotkey func(msg *winapi.Message)) (*UI, error) {

	ui := &UI{actions: actions}
	ui.onHotkey = onHotkey
	hwnd, err := ui.createWindow("Gridda")
	if hwnd == 0 {
		return nil, err
	}
	ui.MainWindow = hwnd

	ui.tray = NewTrayIcon(ui.MainWindow, ui.actions, config)

	ui.createHook()
	return ui, nil
}

// ShowMain will show the main application window
func (ui *UI) ShowMain() {
	winapi.ShowWindow(ui.MainWindow)
}

// Quit will close the message loop
func (ui *UI) Quit() {
	ui.looping = false
	winapi.RemoveHook(ui.hook)
}

// SetGrid will set the active grid to show in the menu
func (ui *UI) SetGrid(config *cmds.Config, index int) {
	ui.tray.SetGrid(config, index)
}

// Shutdown the UI freeing all resources
func (ui *UI) shutdown() {
	ui.tray.Shutdown()
}

// Loop will start the message loop and block until the application is closed
func (ui *UI) Loop() {
	defer ui.shutdown()
	var msg = new(winapi.Message)
	ui.looping = true
	for ui.looping {
		winapi.GetMessage(msg)
		switch msg.Message {
		case winapi.WM_HOTKEY:
			ui.onHotkey(msg)
		case winapi.WM_COMMAND:
			ui.tray.OnCommand(msg)
		default:
			winapi.DispatchMessage(msg)
		}
	}
}

func (ui *UI) createWindow(title string) (winapi.Hwnd, error) {
	return winapi.CreateInactiveWindow(title, 0, 0, 300, 500,
		winapi.NewMessageHandler(winapi.WM_CLOSE, func(hwnd winapi.Hwnd, wParam winapi.Wparam, lParam winapi.LParam) bool {
			fmt.Println("ByeBye world")
			return true
		}),
	)
}

func (ui *UI) createHook() {
	msgID := winapi.RegisterWindowMessage(hookyMsgID)
	winapi.AddMessageHandler(ui.MainWindow, winapi.NewMessageHandler(msgID, func(hwnd winapi.Hwnd, wParam winapi.Wparam, lParam winapi.LParam) bool {
		switch lParam {
		case winapi.HCBT_ACTIVATE:
			ui.actions.OnActivate(winapi.Hwnd(wParam))
		case winapi.HCBT_SETFOCUS:
			ui.actions.OnFocus(winapi.Hwnd(wParam))
		}
		fmt.Printf("Hello Hook! HWND:%d, code:%d\n", wParam, lParam)
		return true
	}))
	ui.hook = winapi.AddCbtHook()
}
