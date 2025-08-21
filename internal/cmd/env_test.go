package cmd

import (
	"testing"
)

func TestEnv(t *testing.T) {
	env := NewEnv()
	env.Set("FOO", "bar")

	out := env.Serialize()

	env2, err := LoadEnv(out)
	if err != nil {
		t.Error("parse error", err)
	}

	if env2.Get("FOO") != "bar" {
		t.Error("FOO != bar", env2.Get("FOO"))
	}

	if env2.Len() != 1 {
		t.Error("len != 1", env2.Len())
	}
}
