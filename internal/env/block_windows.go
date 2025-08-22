package env

// Represents a set of environment variables native to this OS.
type Block = WindowsBlock

// FromMap builds a Block from a map.
func FromMap(m map[string]string) *Block {
	var b Block
	b.fromMap(m)
	return &b
}
