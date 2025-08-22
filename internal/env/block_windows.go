package env

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Represents a set of environment variables native to this OS.
type Block = WindowsBlock

// FromMap builds a Block from a map.
func FromMap(m map[string]string) *Block {
	var b Block
	b.fromMap(m)
	return &b
}

// compareOrdinalIgnoreCase returns -1, 0, +1 using Windows CompareStringOrdinal with bIgnoreCase=TRUE.
func compareOrdinalIgnoreCase(a, b string) int {
	var wa, wb []uint16
	wa, _ = windows.UTF16FromString(a)
	wb, _ = windows.UTF16FromString(b)
	la := len(wa) - 1 // exclude NUL
	lb := len(wb) - 1

	r, _, _ := procCompareStringOrdinal.Call(
		uintptr(unsafe.Pointer(&wa[0])),
		uintptr(la),
		uintptr(unsafe.Pointer(&wb[0])),
		uintptr(lb),
		uintptr(1), // bIgnoreCase = TRUE
	)
	switch r {
	case cstrLessThan:
		return -1
	case cstrEqual:
		return 0
	case cstrGreaterThan:
		return 1
	default:
		// Degenerate fallback: compare lengths then bytes (stable, deterministic)
		if la < lb {
			return -1
		}
		if la > lb {
			return 1
		}
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
}

var (
	modKernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procCompareStringOrdinal = modKernel32.NewProc("CompareStringOrdinal")

	cstrLessThan    uintptr = 1
	cstrEqual       uintptr = 2
	cstrGreaterThan uintptr = 3
)
