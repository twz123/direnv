package env_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/direnv/direnv/v2/internal/env"
)

func TestBlock_Marshal_New(t *testing.T) {
	b, err := json.Marshal(new(env.WindowsBlock))
	if err != nil {
		t.Fatal(err)
	}
	if expected, actual := []byte("{}"), b; !bytes.Equal(expected, actual) {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
