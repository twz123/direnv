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
	vars []entry // kept sorted by key (ignore-case comparator)
}

type entry struct{ name, value string }

func (b *WindowsBlock) fromMap(m map[string]string) {
	b.vars = make([]entry, 0, len(m))
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

// Set assigns a value to a key in the environment.
// The stored key casing is whatever was passed in for the first time.
func (b *WindowsBlock) Set(name, value string) {
	i, found := b.search(name)
	if found {
		// update in place; preserve casing
		b.vars[i].value = value
		return
	}
	// insert while keeping the block sorted by ignore-case comparator
	b.vars = append(b.vars, entry{}) // grow
	copy(b.vars[i+1:], b.vars[i:])
	b.vars[i] = entry{name: name, value: value}
}

// Get retrieves a value from the environment. Returns "" if unset.
func (b *WindowsBlock) Get(name string) string {
	v, _ := b.Lookup(name)
	return v
}

// Lookup retrieves the value for a key and whether it was present.
func (b *WindowsBlock) Lookup(name string) (string, bool) {
	i, found := b.search(name)
	if !found {
		return "", false
	}
	return b.vars[i].value, true
}

// Delete removes a key from the environment (if present).
func (b *WindowsBlock) Delete(name string) {
	if i, found := b.search(name); found {
		copy(b.vars[i:], b.vars[i+1:])
		b.vars = b.vars[:len(b.vars)-1]
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
		for _, e := range b.vars {
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
	return sort.Find(len(b.vars), func(i int) int {
		return compareOrdinalIgnoreCase(name, b.vars[i].name)
	})
}
