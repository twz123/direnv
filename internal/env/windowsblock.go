package env

import (
	"encoding/json"
	"iter"
	"slices"
	"sort"
)

// A Windows-like environment block: case-insensitive (by comparator),
// case-preserving, and kept sorted by name using compareOrdinalIgnoreCase.
type WindowsBlock struct {
	vars []variable // kept sorted by name (ignore-case comparator)
}

type variable struct{ name, value string }

func (b *WindowsBlock) fromMap(m map[string]string) {
	b.vars = make([]variable, 0, len(m))
	for name, value := range m {
		b.Set(name, value)
	}
}

// Copy returns a fresh copy of the block. Because the block is a map under the
// hood, we want to get a copy whenever we mutate it and want to keep the
// original around.
func (b *WindowsBlock) Copy() *WindowsBlock {
	if b == nil {
		return nil
	}
	return &WindowsBlock{slices.Clone(b.vars)}
}

// Sets the value of the environment variable with the given name.
// The stored key casing is whatever was passed in for the first time.
func (b *WindowsBlock) Set(name, value string) {
	if i, found := b.search(name); found {
		// update in place; preserve casing
		b.vars[i].value = value
		return
	} else {
		b.vars = slices.Insert(b.vars, i, variable{name, value})
	}
}

// Retrieves the value of the environment variable with the given name. Returns
// the value, which will be empty if the variable is not present. To distinguish
// between an empty value and an unset value, use [WindowsBlock.Lookup].
func (b *WindowsBlock) Get(name string) string {
	value, _ := b.Lookup(name)
	return value
}

// Retrieves the value of the environment variable with the given name. If the
// variable is present in this block the value (which may be empty) is
// returned and the boolean is true. Otherwise the returned value will be empty
// and the boolean will be false.
func (b *WindowsBlock) Lookup(name string) (string, bool) {
	if b != nil {
		i, found := b.search(name)
		if found {
			return b.vars[i].value, true
		}
	}

	return "", false
}

// Unsets a single environment variable.
func (b *WindowsBlock) Unset(name string) {
	if b == nil {
		return
	}

	if i, found := b.search(name); found {
		b.vars = slices.Delete(b.vars, i, i+1)
	}
}

// Len returns the number of variables in this block.
func (b *WindowsBlock) Len() int {
	if b == nil {
		return 0
	}

	return len(b.vars)
}

// Iterates over all variables in this block.
func (b *WindowsBlock) All() iter.Seq2[string, string] {
	var vars []variable
	if b != nil {
		vars = b.vars
	}

	return func(yield func(string, string) bool) {
		for _, e := range vars {
			if !yield(e.name, e.value) {
				return
			}
		}
	}
}

func (b *WindowsBlock) MarshalJSON() ([]byte, error) {
	block := make(map[string]string, len(b.vars))
	for _, entry := range b.vars {
		block[entry.name] = entry.value
	}
	return json.Marshal(block)
}

func (b *WindowsBlock) UnmarshalJSON(bytes []byte) error {
	m := make(map[string]string)
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}

	b.fromMap(m)
	return nil
}

// search performs a binary search using compareOrdinalIgnoreCase.
// It returns the index where name is/should be, and whether an equal key was found.
func (b *WindowsBlock) search(name string) (idx int, found bool) {
	if b == nil { // FIXME fixup: Implement Windows case-insensitivity
		return
	}
	return sort.Find(len(b.vars), func(i int) int {
		return compareOrdinalIgnoreCase(name, b.vars[i].name)
	})
}
