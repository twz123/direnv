package cmd

import (
	"testing"

	"github.com/direnv/direnv/v2/gzenv"
)

func TestEnv(t *testing.T) {
	env := Env{"FOO": "bar"}

	out := gzenv.Marshal(env)

	var env2 Env
	err := gzenv.Unmarshal(out, &env2)
	if err != nil {
		t.Error("parse error", err)
	}

	if foo := env2["FOO"]; foo != "bar" {
		t.Error("FOO != bar", foo)
	}

	if len(env2) != 1 {
		t.Error("len != 1", len(env2))
	}
}
