package winapi

import (
	"syscall"
	"unsafe"
)

const (
	IMAGE_ICON = 1

	LR_LOADFROMFILE = 0x00000010
	LR_DEFAULTSIZE  = 0x00000040
)

type Icon interface {
	handle() uintptr
}

type icon struct {
	imageHandle uintptr
}

// NewIcon create a new icon from the given file
func NewIcon(fileName string) (Icon, error) {
	var err error
	icon := new(icon)
	icon.imageHandle, _, err = loadImageA.Call(uintptr(0), uintptr(unsafe.Pointer(syscall.StringBytePtr(fileName))), uintptr(IMAGE_ICON), uintptr(0), uintptr(0), uintptr(LR_LOADFROMFILE))
	return icon, ifError(icon.imageHandle == 0, err, "NewIcon")
}

func (icon *icon) handle() uintptr {
	return icon.imageHandle
}
