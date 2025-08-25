package env

import (
	"maps"
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestBasicSetGetLookupDelete(t *testing.T) {
	var b WindowsBlock

	// Initially empty
	if _, ok := b.Lookup("Path"); ok {
		t.Fatalf("unexpected variable present")
	}
	if got := b.Get("Path"); got != "" {
		t.Fatalf("Get returned %q, want empty", got)
	}

	// Set & Lookup
	b.Set("Path", `C:\Windows`)
	if v, ok := b.Lookup("PATH"); !ok || v != `C:\Windows` {
		t.Fatalf("Lookup PATH got (%q, %v), want (%q, true)", v, ok, `C:\Windows`)
	}
	if v := b.Get("path"); v != `C:\Windows` {
		t.Fatalf("Get path got %q", v)
	}

	// Update existing with different casing
	b.Set("pAtH", `C:\Windows;C:\Windows\System32`)
	if v, ok := b.Lookup("PATH"); !ok || v != `C:\Windows;C:\Windows\System32` {
		t.Fatalf("Lookup after update got (%q, %v)", v, ok)
	}

	// FIXME check casing

	// Delete (case-insensitive)
	b.Unset("PATH")
	if _, ok := b.Lookup("Path"); ok {
		t.Fatalf("variable should be gone")
	}
}

func TestCasePreservationOnFirstSet(t *testing.T) {
	var b WindowsBlock
	b.Set("Path", "A")
	b.Set("PATH", "B")
	b.Set("pAtH", "C")

	keys := slices.Collect(maps.Keys(maps.Collect(b.All())))
	if len(keys) != 1 {
		t.Fatalf("want 1 key, got %d: %v", len(keys), keys)
	}
	if keys[0] != "Path" {
		t.Fatalf("didn't preserv key casing = %q, want %q", keys[0], "Path")
	}
	if v, ok := b.Lookup("path"); !ok || v != "C" {
		t.Fatalf("value mismatch: got (%q,%v), want (C,true)", v, ok)
	}
}

func TestOrderingAndKeys(t *testing.T) {
	var b WindowsBlock
	b.Set("TMP", "1")
	b.Set("ComSpec", "2")
	b.Set("Path", "3")
	b.Set("LOCALAPPDATA", "4")

	var got []string
	for name := range b.All() {
		got = append(got, name)
	}
	want := []string{"ComSpec", "LOCALAPPDATA", "Path", "TMP"} // Windows-style ignore-case order
	if !slices.Equal(got, want) {
		t.Fatalf("Keys order = %v, want %v", got, want)
	}

	// Inserting something that sorts to the middle keeps order
	b.Set("PATHeXt", "x")
	got = got[0:0]
	for name := range b.All() {
		got = append(got, name)
	}
	want = []string{"ComSpec", "LOCALAPPDATA", "Path", "PATHeXt", "TMP"}
	if !slices.Equal(got, want) {
		t.Fatalf("Keys order after insert = %v, want %v", got, want)
	}
}

func TestNoDuplicatesUnderDifferentCasing(t *testing.T) {
	var b WindowsBlock
	names := []string{"Var", "VAR", "vAr", "vaR", "var"}
	for i, n := range names {
		b.Set(n, strings.Repeat("x", i+1))
	}

	keys := slices.Collect(maps.Keys(maps.Collect(b.All())))
	if len(keys) != 1 {
		t.Fatalf("expected 1 canonical key, got %d (%v)", len(keys), keys)
	}
	v, ok := b.Lookup("VAR")
	if !ok {
		t.Fatalf("Lookup failed")
	}
	if v != strings.Repeat("x", len(names)) {
		t.Fatalf("unexpected final value: %q", v)
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]string{
		"Path":    "A",
		"TEMP":    "B",
		"ComSpec": "C",
	}
	b := FromMap(m)

	// Verify contents via lookups (case-insensitive)
	for k, want := range m {
		if v, ok := b.Lookup(strings.ToUpper(k)); !ok || v != want {
			t.Fatalf("Lookup %q got (%q,%v), want (%q,true)", strings.ToUpper(k), v, ok, want)
		}
	}
	// Verify order is sorted ignore-case
	var got []string
	for name := range b.All() {
		got = append(got, name)
	}
	want := []string{"ComSpec", "Path", "TEMP"}
	if !slices.Equal(got, want) {
		t.Fatalf("Keys order = %v, want %v", got, want)
	}
}

func TestAsMapSnapshot(t *testing.T) {
	var b WindowsBlock
	b.Set("A", "1")
	b.Set("B", "2")
	m := maps.Collect(b.All())
	if len(m) != 2 || m["A"] != "1" || m["B"] != "2" {
		t.Fatalf("unexpected snapshot: %#v", m)
	}
	// Mutate Block; snapshot must not change.
	b.Set("C", "3")
	delete(m, "A")
	if _, ok := m["C"]; ok {
		t.Fatalf("snapshot should not reflect later changes")
	}
}

func TestUnicodeSanityIgnoreCase(t *testing.T) {
	var b WindowsBlock
	// Greek pi: π (U+03C0), Π (U+03A0)
	b.Set("π", "pi")
	if v, ok := b.Lookup("Π"); !ok || v != "pi" {
		t.Fatalf("Unicode case-insensitive match failed: got (%q,%v)", v, ok)
	}
}

// -------- Fuzzing --------

// Fuzz that repeated Set with different casings yields exactly one key,
// and Lookup works for any casing.
func FuzzSetLookupUniqueness(f *testing.F) {
	seed := []string{"Path", "TMP", "ComSpec", "LOCALAPPDATA"}
	for _, s := range seed {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, base string) {
		if base == "" || strings.ContainsRune(base, '=') {
			t.Skip()
			return
		}
		var b WindowsBlock
		r := rand.New(rand.NewSource(1))
		n := 200
		for i := 0; i < n; i++ {
			// Randomize casing for base
			var sb strings.Builder
			for _, r0 := range base {
				if r.Intn(2) == 0 {
					sb.WriteRune(unicodeToLower(r0))
				} else {
					sb.WriteRune(unicodeToUpper(r0))
				}
			}
			b.Set(sb.String(), strconvI(i))
		}
		// Exactly one canonical key exists
		keys := slices.Collect(maps.Keys(maps.Collect(b.All())))
		if len(keys) != 1 {
			t.Fatalf("expected exactly one canonicalized key for %q, got %d (%v)", base, len(keys), keys)
		}
		// Look up with upper/lower
		if _, ok := b.Lookup(strings.ToUpper(base)); !ok {
			t.Fatalf("Lookup upper failed for %q", base)
		}
		if _, ok := b.Lookup(strings.ToLower(base)); !ok {
			t.Fatalf("Lookup lower failed for %q", base)
		}
	})
}

// Fuzz ordering: inserting many keys in random order must produce Keys()
// that are sorted according to compareOrdinalIgnoreCase.
// func FuzzOrdering(f *testing.F) {
// 	seeds := [][]string{
// 		{"a", "B", "aa", "ZZ", "z", "Path", "pathext", "PATHeXt", "ComSpec"},
// 		{"π", "Π", "Τ", "τ", "TMP", "tmp"},
// 	}
// 	for _, s := range seeds {
// 		f.Add(s)
// 	}
// 	f.Fuzz(func(t *testing.T, inputs []string) {
// 		var b WindowsBlock
// 		for _, k := range inputs {
// 			if k == "" || strings.ContainsRune(k, '=') {
// 				continue
// 			}
// 			b.Set(k, "x")
// 		}
// 		got := slices.Collect(maps.Keys(maps.Collect(b.All())))
// 		sorted := slices.Clone(got)
// 		sort.Slice(sorted, func(i, j int) bool {
// 			return compareOrdinalIgnoreCase(sorted[i], sorted[j]) < 0
// 		})
// 		if !slices.Equal(got, sorted) {
// 			t.Fatalf("Keys not in comparator order.\n got: %v\nwant: %v", got, sorted)
// 		}
// 	})
// }

// ---- tiny helpers used by fuzz ----

func unicodeToLower(r rune) rune { return []rune(strings.ToLower(string(r)))[0] }
func unicodeToUpper(r rune) rune { return []rune(strings.ToUpper(string(r)))[0] }
func strconvI(i int) string      { return strconvFormatInt(int64(i)) }

// Avoid importing strconv just for Int->string in fuzz
func strconvFormatInt(i int64) string {
	return fmtInt(i)
}

// minimal int->string for non-negative i used above
func fmtInt(i int64) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + (i % 10))
		i /= 10
	}
	return string(buf[pos:])
}
