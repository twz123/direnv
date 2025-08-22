package env

// Block is a map representation of environment variables.
type Block map[string]string

// Copy returns a fresh copy of the env. Because the env is a map under the
// hood, we want to get a copy whenever we mutate it and want to keep the
// original around.
func (env Block) Copy() Block {
	newEnv := make(Block)

	for key, value := range env {
		newEnv[key] = value
	}

	return newEnv
}
