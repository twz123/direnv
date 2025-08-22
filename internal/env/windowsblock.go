package env

import (
	"encoding/json"
	"iter"
	"maps"
)

// Represents a set of environment variables on Windows operating systems.
//
// The underlying representation is hidden so that custom implementations on
// how to store and retrieve values can be provided later on.
type WindowsBlock struct {
	vars map[string]string
}

// Copy returns a fresh copy of the block. Because the block is a map under the
// hood, we want to get a copy whenever we mutate it and want to keep the
// original around.
func (b *WindowsBlock) Copy() *WindowsBlock {
	if b == nil {
		return nil
	}
	return &WindowsBlock{maps.Clone(b.vars)}
}

// Set assigns a value to a key in the environment.
func (b *WindowsBlock) Set(key, value string) {
	if b.vars == nil {
		b.vars = map[string]string{key: value}
	} else {
		b.vars[key] = value
	}
}

// Get retrieves a value from the environment. Returns "" if unset.
func (b *WindowsBlock) Get(key string) (val string) {
	if b != nil {
		val = b.vars[key]
	}
	return
}

// Lookup retrieves the value for a key and whether it was present.
func (b *WindowsBlock) Lookup(key string) (val string, ok bool) {
	if b != nil {
		val, ok = b.vars[key]

	}
	return
}

// Delete removes a key from the environment.
func (b *WindowsBlock) Delete(key string) {
	if b != nil {
		delete(b.vars, key)
	}
}

// Len returns the number of entries in the environment.
func (b *WindowsBlock) Len() int {
	if b == nil {
		return 0
	}

	return len(b.vars)
}

// All returns an iterator over all key/value pairs in the environment.
func (b *WindowsBlock) All() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for k, v := range b.vars {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (b *WindowsBlock) MarshalJSON() ([]byte, error) {
	if len(b.vars) < 1 {
		return []byte("{}"), nil
	}

	return json.Marshal(b.vars)
}

func (b *WindowsBlock) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &b.vars)
}
