package ui

import (
	"fmt"

	"sonnenfroh.de/test/cmds"
	"sonnenfroh.de/test/winapi"
)

// UI is the base instance for all ui related tasks
type UI struct {
	MainWindow winapi.Hwnd
	tray       TrayIcon
	looping    bool
	onHotkey   func(msg *winapi.Message)
	actions    Actions
}

// Actions describe all actions which can be triggered by the UI
type Actions interface {
	TrayIconActions
}

// New creates a new ui instance
func New(actions Actions, config *cmds.Config, onHotkey func(msg *winapi.Message)) (*UI, error) {
	ui := &UI{actions: actions}
	ui.onHotkey = onHotkey
	hwnd, err := ui.createWindow("Grida")
	if hwnd == 0 {
		return nil, err
	}
	ui.MainWindow = hwnd

	// ui.createMenu()

	ui.tray = NewTrayIcon(ui.MainWindow, ui.actions, config)

	return ui, nil
}

// ShowMain will show the main application window
func (ui *UI) ShowMain() {
	winapi.ShowWindow(ui.MainWindow)
}

// Quit will close the message loop
func (ui *UI) Quit() {
	ui.looping = false
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
