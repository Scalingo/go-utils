package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitMapFromEnv(t *testing.T) {
	defaults := map[string]string{
		"VAR1": "val1",
		"VAR2": "val2",
	}
	os.Setenv("VAR2", "pipomolo")

	res := InitMapFromEnv(defaults)
	assert.Equal(t, res["VAR1"], "val1")
	assert.Equal(t, res["VAR2"], "pipomolo")
}
