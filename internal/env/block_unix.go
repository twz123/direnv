//go:build !windows

package env

import "maps"

// Represents a set of environment variables native to this OS.
type Block = UNIXBlock

// FromMap builds a Block from a map.
func FromMap(m map[string]string) *Block {
	return &Block{maps.Clone(m)}
}
