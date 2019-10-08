package winapi

import "fmt"

// Hwnd defines a handle for a window
type Hwnd uintptr
type Wparam uintptr
type LParam uintptr

// Point defines a position of e.g. a window
type Point struct {
	x int32
	y int32
}

// Rect defines a rectangle by coordinates of upper-left and lower-right corner
type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// Width returns the width of the rectangle
func (rect *Rect) Width() int32 {
	return rect.Right - rect.Left
}

// Height returns the width of the rectangle
func (rect *Rect) Height() int32 {
	return rect.Bottom - rect.Top
}

func (rect *Rect) String() string {
	return fmt.Sprintf("Rect: top:%d, left:%d, bottom:%d, right:%d\n", rect.Top, rect.Left, rect.Bottom, rect.Right)
}

// Message defines a windows message queue event
type Message struct {
	Hwnd    Hwnd
	Message uint
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Point   Point
}
