package cmd

import (
	"encoding/json"
	"iter"
	"os"
	"strings"

	"github.com/direnv/direnv/v2/gzenv"
)

// Env represents a set of environment variables.
//
// The underlying representation is hidden so that custom implementations on
// how to store and retrieve values can be provided later on.
type Env struct {
	vars map[string]string
}

// GetEnv turns the classic unix environment variables into a map of
// key->values which is more handy to work with.
//
// NOTE:  We don't support having two variables with the same name.
// I've never seen it used in the wild but according to POSIX it's allowed.
func GetEnv() Env {
	env := NewEnv()

	for _, kv := range os.Environ() {
		kv2 := strings.SplitN(kv, "=", 2)
		key := kv2[0]
		value := kv2[1]
		env.Set(key, value)
	}

	return env
}

// CleanContext removes all the direnv-related environment variables. Call
// this after reverting the environment, otherwise direnv will just be amnesic
// about the previously-loaded environment.
func (env Env) CleanContext() {
	env.Delete(DIRENV_DIFF)
	env.Delete(DIRENV_DIR)
	env.Delete(DIRENV_FILE)
	env.Delete(DIRENV_DUMP_FILE_PATH)
	env.Delete(DIRENV_WATCHES)
}

// LoadEnv unmarshals the env back from a gzenv string
func LoadEnv(gzenvStr string) (env Env, err error) {
	m := map[string]string{}
	if err = gzenv.Unmarshal(gzenvStr, &m); err != nil {
		return Env{}, err
	}
	env = NewEnv()
	for k, v := range m {
		env.Set(k, v)
	}
	return env, nil
}

// LoadEnvJSON unmarshals the env back from a JSON string
func LoadEnvJSON(jsonBytes []byte) (env Env, err error) {
	m := map[string]string{}
	if err = json.Unmarshal(jsonBytes, &m); err != nil {
		return Env{}, err
	}
	env = NewEnv()
	for k, v := range m {
		env.Set(k, v)
	}
	return env, nil
}

// Copy returns a fresh copy of the env. Because the env is a map under the
// hood, we want to get a copy whenever we mutate it and want to keep the
// original around.
func (env Env) Copy() Env {
	newEnv := NewEnv()
	for key, value := range env.All() {
		newEnv.Set(key, value)
	}
	return newEnv
}

// ToGoEnv should really be named ToUnixEnv. It turns the env back into a list
// of "key=value" strings like returns by os.Environ().
func (env Env) ToGoEnv() []string {
	goEnv := make([]string, env.Len())
	index := 0
	for key, value := range env.All() {
		goEnv[index] = strings.Join([]string{key, value}, "=")
		index++
	}
	return goEnv
}

// ToShell outputs the environment into an evaluatable string that is
// understood by the target shell
func (env Env) ToShell(shell Shell) (string, error) {
	e := make(ShellExport)
	for key, value := range env.All() {
		e.Add(key, value)
	}
	return shell.Export(e)
}

// Serialize marshals the env into the gzenv format
func (env Env) Serialize() string {
	return gzenv.Marshal(env.vars)
}

// Diff returns the diff between the current env and the passed env
func (env Env) Diff(other Env) *EnvDiff {
	return BuildEnvDiff(env, other)
}

// Fetch tries to get the value associated with the given 'key', or returns
// the provided default if none is set.
//
// Note that empty environment variables are considered to be set.
func (env Env) Fetch(key, def string) string {
	v, ok := env.Lookup(key)
	if !ok {
		v = def
	}
	return v
}

// NewEnv returns an empty environment.
func NewEnv() Env {
	return Env{vars: make(map[string]string)}
}

// NewEnvFrom returns an environment initialized with the provided map.
func NewEnvFrom(m map[string]string) Env {
	env := NewEnv()
	for k, v := range m {
		env.Set(k, v)
	}
	return env
}

// Set assigns a value to a key in the environment.
func (env Env) Set(key, value string) {
	env.vars[key] = value
}

// Get retrieves a value from the environment. Returns "" if unset.
func (env Env) Get(key string) string {
	return env.vars[key]
}

// Lookup retrieves the value for a key and whether it was present.
func (env Env) Lookup(key string) (string, bool) {
	v, ok := env.vars[key]
	return v, ok
}

// Delete removes a key from the environment.
func (env Env) Delete(key string) {
	delete(env.vars, key)
}

// Len returns the number of entries in the environment.
func (env Env) Len() int {
	return len(env.vars)
}

// All returns an iterator over all key/value pairs in the environment.
func (env Env) All() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for k, v := range env.vars {
			if !yield(k, v) {
				return
			}
		}
	}
}
