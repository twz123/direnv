//go:build !windows

package env

import (
	"strings"
)

// Represents a set of environment variables native to this OS.
type Block = WindowsBlock

// FromMap builds a Block from a map.
func FromMap(m map[string]string) *Block {
	var b Block
	b.fromMap(m)
	return &b
}

// // FromMap builds a Block from a map.
// func FromMap(m map[string]string) *Block {
// 	return &Block{maps.Clone(m)}
// }

// compareOrdinalIgnoreCase returns -1, 0, +1 in a way that *simulates*
// Windows' ordinal-ignore-case using Go's standard library.
// Note: This is NOT an exact replica of Windows' case mapping rules; it uses
// Unicode simple case-folding (locale-independent).
func compareOrdinalIgnoreCase(a, b string) int {
	// Fast path: equal by Unicode case-folding
	if strings.EqualFold(a, b) {
		return 0
	}
	// Order by uppercased strings to get a stable total order.
	au := strings.ToUpper(a)
	bu := strings.ToUpper(b)
	switch {
	case au < bu:
		return -1
	case au > bu:
		return 1
	default:
		// If folded forms collide but EqualFold said not equal,
		// fall back to raw compare to break ties deterministically.
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
}
