package env

import (
	"encoding/json"
	"iter"
	"maps"
)

// Represents a set of environment variables on UNIX operating systems.
//
// The underlying representation is hidden so that custom implementations on
// how to store and retrieve values can be provided later on.
type UNIXBlock struct {
	vars map[string]string
}

// Copy returns a fresh copy of the block. Because the block is a map under the
// hood, we want to get a copy whenever we mutate it and want to keep the
// original around.
func (b *UNIXBlock) Copy() *UNIXBlock {
	if b == nil {
		return nil
	}
	return &UNIXBlock{maps.Clone(b.vars)}
}

// Sets the value of the environment variable with the given name.
func (b *UNIXBlock) Set(key, value string) {
	if b.vars == nil {
		b.vars = map[string]string{key: value}
	} else {
		b.vars[key] = value
	}
}

// Retrieves the value of the environment variable with the given name. Returns
// the value, which will be empty if the variable is not present. To distinguish
// between an empty value and an unset value, use [UNIXBlock.Lookup].
func (b *UNIXBlock) Get(name string) string {
	if b != nil {
		return b.vars[name]
	}

	return ""
}

// Retrieves the value of the environment variable with the given name. If the
// variable is present in this block the value (which may be empty) is
// returned and the boolean is true. Otherwise the returned value will be empty
// and the boolean will be false.
func (b *UNIXBlock) Lookup(name string) (string, bool) {
	if b != nil {
		value, found := b.vars[name]
		return value, found
	}
	return "", false
}

// Unsets a single environment variable.
func (b *UNIXBlock) Unset(name string) {
	if b == nil {
		return
	}

	delete(b.vars, name)
}

// Len returns the number of variables in this block.
func (b *UNIXBlock) Len() int {
	if b == nil {
		return 0
	}

	return len(b.vars)
}

// Iterates over all variables in this block.
func (b *UNIXBlock) All() iter.Seq2[string, string] {
	var vars map[string]string
	if b != nil {
		vars = b.vars
	}

	return func(yield func(string, string) bool) {
		for name, value := range vars {
			if !yield(name, value) {
				return
			}
		}
	}
}

func (b *UNIXBlock) MarshalJSON() ([]byte, error) {
	if len(b.vars) < 1 {
		return []byte("{}"), nil
	}

	return json.Marshal(b.vars)
}

func (b *UNIXBlock) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &b.vars)
}
