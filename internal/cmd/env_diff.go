package cmd

import (
	"fmt"
	"strings"

	"github.com/direnv/direnv/v2/gzenv"
)

// IgnoredKeys is list of keys we don't want to deal with
var IgnoredKeys = map[string]bool{
	// direnv env config
	"DIRENV_CONFIG": true,
	"DIRENV_BASH":   true,

	// should only be available inside of the .envrc or .env
	"DIRENV_IN_ENVRC": true,

	"COMP_WORDBREAKS": true, // Avoids segfaults in bash
	"PS1":             true, // PS1 should not be exported, fixes problem in bash

	// variables that should change freely
	"OLDPWD":    true,
	"PWD":       true,
	"SHELL":     true,
	"SHELLOPTS": true,
	"SHLVL":     true,
	"_":         true,
}

// EnvDiff represents the diff between two environments
type EnvDiff struct {
	changes []envDiff
}

// BuildEnvDiff analyses the changes between from and to and builds an
// EnvDiff out of it.
func BuildEnvDiff(from, to Env) *EnvDiff {
	var diff EnvDiff
	for name, fromValue := range from.All() {
		toValue, found := to.Lookup(name)
		if found {
			diff.change(name, fromValue, toValue)
		} else {
			diff.remove(name, fromValue)
		}
	}
	for name, value := range to.All() {
		if _, found := from.Lookup(name); !found {
			diff.add(name, value)
		}
	}

	return &diff
}

// LoadEnvDiff unmarshalls a gzenv string back into an EnvDiff.
func LoadEnvDiff(gzenvStr string) (diff *EnvDiff, err error) {
	diff = new(EnvDiff)
	err = gzenv.Unmarshal(gzenvStr, &diff.changes)
	return
}

type EnvDiffApplier interface {
	Set(name, value string)
	Unset(name string)
}

func (diff *EnvDiff) Apply(applier EnvDiffApplier) {
	for _, change := range diff.changes {
		switch change.typ() {
		case addedToEnv:
			applier.Set(change.name(), change.value1())
		case changedInEnv:
			applier.Set(change.name(), change.value2())
		case removedFromEnv:
			applier.Unset(change.name())
		}
	}
}

func (diff *EnvDiff) Revert(applier EnvDiffApplier) {
	for _, change := range diff.changes {
		switch change.typ() {
		case addedToEnv:
			applier.Unset(change.name())
		case changedInEnv, removedFromEnv:
			applier.Set(change.name(), change.value1())
		}
	}
}

// ToShell applies the env diff as a set of commands that are understood by
// the target `shell`. The outputted string is then meant to be evaluated in
// the target shell.
func (diff *EnvDiff) ToShell(shell Shell) (string, error) {
	e := make(ShellExport)
	diff.Apply(e)
	return shell.Export(e)
}

// Serialize marshalls the environment diff to the gzenv format.
func (diff *EnvDiff) Serialize() string {
	return gzenv.Marshal(diff.changes)
}

//// Utils

// IgnoredEnv returns true if the key should be ignored in environment diffs.
func IgnoredEnv(key string) bool {
	if strings.HasPrefix(key, "__fish") {
		return true
	}
	if strings.HasPrefix(key, "BASH_FUNC_") {
		return true
	}
	_, found := IgnoredKeys[key]
	return found
}

type envDiffType byte

const (
	addedToEnv     envDiffType = '+'
	changedInEnv   envDiffType = '~'
	removedFromEnv envDiffType = '-'
)

type envDiff []string

func (d *EnvDiff) add(name, value string) {
	d.changes = append(d.changes, envDiff([]string{fmt.Sprintf("%c%s", addedToEnv, name), value}))
}

func (d *EnvDiff) change(name, fromValue, toValue string) {
	d.changes = append(d.changes, envDiff([]string{fmt.Sprintf("%c%s", changedInEnv, name), fromValue, toValue}))
}

func (d *EnvDiff) remove(name, oldValue string) {
	d.changes = append(d.changes, envDiff([]string{fmt.Sprintf("%c%s", removedFromEnv, name), oldValue}))
}

func (d envDiff) typ() envDiffType { return envDiffType([]string(d)[0][0]) }
func (d envDiff) name() string     { return []string(d)[0][1:] }
func (d envDiff) value1() string   { return []string(d)[1] }
func (d envDiff) value2() string   { return []string(d)[2] }
