//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import (
	"errors"

	"github.com/moonblue4242/Gridda/winapi"
)

// Grid defines a basic grid on the screen at which windows can be snapped to
type Grid struct {
	Name    string
	Columns []int
	Rows    []int
	Presets []Presets
}

// Presets defines values for windows which are set as defaults
type Presets struct {
	Executable string
	Span       *Span
}

func (grid *Grid) columnsWeight() (weight int) {
	weight = 0
	for _, column := range grid.Columns {
		weight += column
	}
	return
}

func (grid *Grid) rowsWeight() (weight int) {
	weight = 0
	for _, rows := range grid.Rows {
		weight += rows
	}
	return
}

// Validate determines if the grid object is setup correctly
func (grid *Grid) Validate() (err error) {
	err = nil
	if len(grid.Name) < 1 {
		err = errors.New("Grid must have a name with at least one character")
	}
	return
}

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
	target := new(targetWindow)
	target.hwnd = winapi.GetForegroundWindow()
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
