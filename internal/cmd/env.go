package cmd

import (
	"os"
	"strings"

	"github.com/direnv/direnv/v2/internal/env"
)

// Env is a map representation of environment variables.
//
// Deprecated: Use env.Block directly.
type Env = env.Block

// GetEnv turns the classic unix environment variables into a map of
// key->values which is more handy to work with.
//
// NOTE:  We don't support having two variables with the same name.
// I've never seen it used in the wild but according to POSIX it's allowed.
func GetEnv() Env {
	env := make(Env)

	for _, kv := range os.Environ() {
		kv2 := strings.SplitN(kv, "=", 2)

		key := kv2[0]
		value := kv2[1]

		env[key] = value
	}

	return env
}

// CleanContext removes all the direnv-related environment variables. Call
// this after reverting the environment, otherwise direnv will just be amnesic
// about the previously-loaded environment.
func CleanContext(env Env) {
	delete(env, DIRENV_DIFF)
	delete(env, DIRENV_DIR)
	delete(env, DIRENV_FILE)
	delete(env, DIRENV_DUMP_FILE_PATH)
	delete(env, DIRENV_WATCHES)
}

// ToGoEnv should really be named ToUnixEnv. It turns the env back into a list
// of "key=value" strings like returns by os.Environ().
func ToGoEnv(env Env) []string {
	goEnv := make([]string, len(env))
	index := 0
	for key, value := range env {
		goEnv[index] = strings.Join([]string{key, value}, "=")
		index++
	}
	return goEnv
}

// ToShell outputs the environment into an evaluatable string that is
// understood by the target shell
func ToShell(env Env, shell Shell) (string, error) {
	e := make(ShellExport)

	for key, value := range env {
		e.Add(key, value)
	}

	return shell.Export(e)
}
