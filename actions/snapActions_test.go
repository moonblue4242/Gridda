package actions

import (
	"testing"

	"github.com/moonblue4242/Gridda/winapi"
	"github.com/stretchr/testify/assert"
)

var (
	snapActionsMock  = NewSnapActions()
	targetWindowMock = new(TargetWindowMock)
	activeConfigMock = new(ActiveConfigMock)

	moveFunc func(left int, top int, width int, height int)
)

type size struct {
	left   int
	top    int
	width  int
	height int
}

type TargetWindowMock struct {
	hwnd        winapi.Hwnd
	size        *winapi.Rect
	desktopSize *winapi.Rect
	deltaH      int
	deltaV      int
}

type ActiveConfigMock struct {
	grid *Grid
}

func (activeConfigMock *ActiveConfigMock) Grid() *Grid {
	return activeConfigMock.grid
}

func (activeConfigMock *ActiveConfigMock) GridIndex() int {
	return 0
}
func (activeConfigMock *ActiveConfigMock) OnPreviousGrid() {}
func (activeConfigMock *ActiveConfigMock) OnNextGrid()     {}

func (targetWindowMock *TargetWindowMock) Hwnd() winapi.Hwnd {
	return targetWindowMock.hwnd
}

func (targetWindowMock *TargetWindowMock) Size() *winapi.Rect {
	return targetWindowMock.size
}

func (targetWindowMock *TargetWindowMock) DesktopSize() *winapi.Rect {
	return targetWindowMock.desktopSize
}

func (targetWindowMock *TargetWindowMock) Delta() (deltaH int, deltaV int) {
	return deltaH, deltaV
}

func (targetWindowMock *TargetWindowMock) Move(left int, top int, right int, bottom int) {
	moveFunc(left, top, right, bottom)
}

func TestSnapLeft(t *testing.T) {
	targetWindowMock.size = &winapi.Rect{Left: 100, Top: 100, Right: 600, Bottom: 400}
	targetWindowMock.desktopSize = &winapi.Rect{Left: 0, Top: 0, Right: 800, Bottom: 600}

	activeConfigMock.grid = new(Grid)
	activeConfigMock.grid.Columns = []int{1, 2, 1}
	activeConfigMock.grid.Rows = []int{1, 1}
	moveFunc = assertMove(t, &size{0, 0, 200, 300})
	snapActionsMock.snapLeft(targetWindowMock, activeConfigMock)
}

func TestSnapRight(t *testing.T) {
	targetWindowMock.size = &winapi.Rect{Left: 100, Top: 100, Right: 600, Bottom: 400}
	targetWindowMock.desktopSize = &winapi.Rect{Left: 0, Top: 0, Right: 800, Bottom: 600}

	activeConfigMock.grid = new(Grid)
	activeConfigMock.grid.Columns = []int{1, 2, 1}
	activeConfigMock.grid.Rows = []int{1, 1}
	moveFunc = assertMove(t, &size{200, 0, 400, 300})
	snapActionsMock.snapRight(targetWindowMock, activeConfigMock)
}

func TestSnapTop(t *testing.T) {
	targetWindowMock.size = &winapi.Rect{Left: 100, Top: 100, Right: 600, Bottom: 400}
	targetWindowMock.desktopSize = &winapi.Rect{Left: 0, Top: 0, Right: 800, Bottom: 600}

	activeConfigMock.grid = new(Grid)
	activeConfigMock.grid.Columns = []int{1, 2, 1}
	activeConfigMock.grid.Rows = []int{1, 1}
	moveFunc = assertMove(t, &size{0, 0, 200, 300})
	snapActionsMock.snapTop(targetWindowMock, activeConfigMock)
}

func TestSnapBottom(t *testing.T) {
	targetWindowMock.size = &winapi.Rect{Left: 100, Top: 100, Right: 600, Bottom: 400}
	targetWindowMock.desktopSize = &winapi.Rect{Left: 0, Top: 0, Right: 800, Bottom: 600}

	activeConfigMock.grid = new(Grid)
	activeConfigMock.grid.Columns = []int{1, 2, 1}
	activeConfigMock.grid.Rows = []int{1, 1}
	moveFunc = assertMove(t, &size{0, 300, 200, 300})
	snapActionsMock.snapBottom(targetWindowMock, activeConfigMock)
}

func TestSpanHorizontal(t *testing.T) {
	assert := assert.New(t)
	targetWindowMock.hwnd = 4711
	targetWindowMock.size = &winapi.Rect{Left: 100, Top: 100, Right: 600, Bottom: 400}
	targetWindowMock.desktopSize = &winapi.Rect{Left: 0, Top: 0, Right: 800, Bottom: 600}

	activeConfigMock.grid = new(Grid)
	activeConfigMock.grid.Columns = []int{1, 2, 1}
	activeConfigMock.grid.Rows = []int{1, 1}
	moveFunc = assertMove(t, &size{0, 0, 600, 300})
	// execute
	snapActionsMock.spanHorizontal(targetWindowMock, activeConfigMock, true)
	// verify
	item, ok := snapActionsMock.spans[0][targetWindowMock.hwnd]
	assert.True(ok, "Span should exists")
	assert.Equal(2, item.Columns, "Increased column span expected")
	assert.Equal(1, item.Rows, "default for row")
}

func assertMove(t *testing.T, expect *size) func(left int, top int, width int, height int) {
	assert := assert.New(t)
	return func(left int, top int, width int, height int) {
		assert.Equal(expect.left, left, "Left")
		assert.Equal(expect.top, top, "Top")
		assert.Equal(expect.width, width, "Width")
		assert.Equal(expect.height, height, "Height")
	}
}
