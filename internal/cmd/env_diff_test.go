package cmd

import (
	"reflect"
	"testing"

	"github.com/direnv/direnv/v2/internal/env"
)

func TestEnvDiff(t *testing.T) {
	diff := &EnvDiff{map[string]string{"FOO": "bar"}, map[string]string{"BAR": "baz"}}

	out := diff.Serialize()

	diff2, err := LoadEnvDiff(out)
	if err != nil {
		t.Error("parse error", err)
	}

	if len(diff2.Prev) != 1 {
		t.Error("len(diff2.prev) != 1", len(diff2.Prev))
	}

	if len(diff2.Next) != 1 {
		t.Error("len(diff2.next) != 0", len(diff2.Next))
	}
}

// Issue #114
// Check that empty environment variables correctly appear in the diff
func TestEnvDiffEmptyValue(t *testing.T) {
	var before, after env.Block
	after.Set("FOO", "")

	diff := BuildEnvDiff(&before, &after)

	expected := map[string]string{"FOO": ""}
	if !reflect.DeepEqual(diff.Next, expected) {
		t.Errorf("diff.Next != after (%#+v != %#+v)", diff.Next, expected)
	}
}

func TestIgnoredEnv(t *testing.T) {
	if !IgnoredEnv(DIRENV_BASH) {
		t.Fail()
	}
	if IgnoredEnv(DIRENV_DIFF) {
		t.Fail()
	}
	if !IgnoredEnv("_") {
		t.Fail()
	}
	if !IgnoredEnv("__fish_foo") {
		t.Fail()
	}
	if !IgnoredEnv("__fishx") {
		t.Fail()
	}
}
