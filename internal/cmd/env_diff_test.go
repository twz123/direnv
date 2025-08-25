package cmd

import (
	"reflect"
	"testing"

	"github.com/direnv/direnv/v2/internal/env"
)

func TestEnvDiff(t *testing.T) {
	var e1, e2 env.WindowsBlock
	e1.Set("FOO", "bar")
	e2.Set("BAR", "upper")

	diff := BuildEnvDiff(&e1, &e2)

	out := diff.Serialize()

	diff2, err := LoadEnvDiff(out)
	if err != nil {
		t.Error("parse error", err)
	}

	want := &EnvDiff{[]envDiff{
		{"-FOO", "bar"},
		{"+BAR", "upper"},
	}}
	if !reflect.DeepEqual(diff2, want) {
		t.Errorf("Unexpected diff: want %#+v, got %#+v", want, diff2)
	}
}

func TestBuildEnvDiff_Windows(t *testing.T) {
	var e1, e2 env.WindowsBlock
	e1.Set("Path", "orig")
	e2.Set("PATH", "upper")

	diff := BuildEnvDiff(&e1, &e2)

	out := diff.Serialize()

	diff2, err := LoadEnvDiff(out)
	if err != nil {
		t.Error("parse error", err)
	}

	status := diffStatus(diff2)
	if expected := "~Path"; expected != status {
		t.Errorf("diffStatus(diff2) != %q %q", expected, status)
	}
}

// Issue #114
// Check that empty environment variables correctly appear in the diff
func TestEnvDiffEmptyValue(t *testing.T) {
	var before, after env.Block
	after.Set("FOO", "")

	diff := BuildEnvDiff(&before, &after)

	want := &EnvDiff{[]envDiff{{"+FOO", ""}}}
	if !reflect.DeepEqual(diff, want) {
		t.Errorf("Unexpected diff: want %#+v, got %#+v", want, diff)
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
