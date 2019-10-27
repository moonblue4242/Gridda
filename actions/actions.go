//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import (
	"strings"

	"github.com/moonblue4242/Gridda/winapi"
)

// ActiveConfig defines the currently active configuration along with additional data
// like row or column span
type ActiveConfig interface {
	Grid() *Grid
	GridIndex() int
	OnPreviousGrid()
	OnNextGrid()
}

// ActionHandler defines the handler functions an actions must comply with
type ActionHandler func(activeConfig ActiveConfig)

func MoveRight(activeConfig ActiveConfig) {
	target := GetTarget()
	target.Move(int(target.Size().Left+50), int(target.Size().Top), int(target.Size().Width()), int(target.Size().Height()))
}

func MoveLeft(activeConfig ActiveConfig) {
	target := GetTarget()
	target.Move(int(target.Size().Left-50), int(target.Size().Top), int(target.Size().Width()), int(target.Size().Height()))
}

func ToPreviousGrid(activeConfig ActiveConfig) {
	activeConfig.OnPreviousGrid()
}

func ToNextGrid(activeConfig ActiveConfig) {
	activeConfig.OnNextGrid()
}

// TargetWindow provides data and functions specific to the window along with aditional general informations
// like size of the desktop
type TargetWindow interface {
	Hwnd() winapi.Hwnd
	Size() *winapi.Rect
	DesktopSize() *winapi.Rect
	Delta() (deltaH int, deltaV int)
	Move(left int, top int, width int, height int)
	ModuleName() (name string, err error)
}

type targetWindow struct {
	hwnd        winapi.Hwnd
	size        winapi.Rect
	desktopSize *winapi.Rect
	deltaH      int
	deltaV      int
}

// GetTarget will retrieve the currently focus window
func GetTarget() TargetWindow {
	return GetTargetFromHandle(winapi.GetForegroundWindow())
}

// GetTargetFromHandle will create a target window out of the given windows window handle
func GetTargetFromHandle(hwnd winapi.Hwnd) TargetWindow {
	target := new(targetWindow)
	target.hwnd = hwnd
	winapi.GetWindowRect(target.hwnd, &target.size)
	target.desktopSize, _ = winapi.GetWorkArea()
	target.updateBorderCorrections()
	return target
}

// Hwnd returns the unique window handle
func (target *targetWindow) Hwnd() winapi.Hwnd {
	return target.hwnd
}

// Size returns the size of the target window
func (target *targetWindow) Size() *winapi.Rect {
	return &target.size
}

// DesktopSize returns the size of the work area (desktop)
func (target *targetWindow) DesktopSize() *winapi.Rect {
	return target.desktopSize
}

// Move will move and size the given window, taking into acount the deltas provided
func (target *targetWindow) Move(left int, top int, width int, height int) {
	// winapi.MoveWindow(target.hwnd, left-target.deltaH, top, width+target.deltaH*2, height+target.deltaV)
	winapi.MoveWindow(target.hwnd, left, top, width, height)
}

// Delta returns the border delta of the given window
func (target *targetWindow) Delta() (deltaH int, deltaV int) {
	return target.deltaH, target.deltaV
}

// calcBorderCorrectiosn will calculate corrections regaring the shadowed borders of Win 10
func (target *targetWindow) updateBorderCorrections() {
	clientArea, _ := winapi.BorderEx(target.hwnd)
	target.deltaH = int(target.size.Width()-clientArea.Width()) / 2
	target.deltaV = int(target.size.Height() - clientArea.Height())
	// correct height off by one
	if target.deltaV > 0 {
		target.deltaV++
	}
	return
}

// moduleName will return the name of the module the window is associated with
func (target *targetWindow) ModuleName() (name string, err error) {
	processID := winapi.GetWindowThreadProcessID(target.Hwnd())
	if processHandle, err := winapi.OpenProcess(processID); err == nil {
		name = strings.ToUpper(winapi.GetModuleBaseName(processHandle))
	}
	return
}
