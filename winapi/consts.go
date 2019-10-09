package winapi

const (
	HWND_TOP = 0

	MB_YESNOCANCEL = 0x00000003

	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
	MOD_NOREPEAT = 0x4000

	SPI_GETWORKAREA = 0x0030

	VK_LEFT  = 0x25
	VK_UP    = 0x26
	VK_RIGHT = 0x27
	VK_DOWN  = 0x28
	VK_PRIOR = 0x21
	VK_NEXT  = 0x22

	VK_A = 0x41
	VK_Z = 0x5A

	WM_QUIT          = 0x0012
	WM_CLOSE         = 0x0010
	WM_DESTROY       = 0x0002
	WM_HOTKEY        = 0x0312
	WM_LBUTTONDBLCLK = 0x203
	WM_RBUTTONUP     = 0x0205
	WM_RBUTTONDOWN   = 0x0204
)