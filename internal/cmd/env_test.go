package cmd

import (
	"testing"

	"github.com/direnv/direnv/v2/gzenv"
	"github.com/direnv/direnv/v2/internal/env"
)

func TestEnv(t *testing.T) {
	env := new(env.Block)
	env.Set("FOO", "bar")

	out := gzenv.Marshal(env)

	var env2 Env
	err := gzenv.Unmarshal(out, &env2)
	if err != nil {
		t.Error("parse error", err)
	}

	if foo := env2.Get("FOO"); foo != "bar" {
		t.Error("FOO != bar", foo)
	}

	if env2.Len() != 1 {
		t.Error("len != 1", env2.Len())
	}
}
