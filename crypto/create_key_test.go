package crypto

import (
	"crypto/hmac"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateKey(t *testing.T) {
	t.Parallel()

	key, err := CreateKey(32)
	require.NoError(t, err)
	assert.Len(t, key, 32)

	key2, err := CreateKey(32)
	require.NoError(t, err)
	assert.Len(t, key2, 32)

	require.False(t, hmac.Equal(key, key2))
}
