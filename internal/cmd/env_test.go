package cmd

import (
	"testing"

	"github.com/direnv/direnv/v2/gzenv"
)

func TestEnv(t *testing.T) {
	env := Env{"FOO": "bar"}

	out := gzenv.Marshal(env)

	env2, err := LoadEnv(out)
	if err != nil {
		t.Error("parse error", err)
	}

	if env2["FOO"] != "bar" {
		t.Error("FOO != bar", env2["FOO"])
	}

	if len(env2) != 1 {
		t.Error("len != 1", len(env2))
	}
}
