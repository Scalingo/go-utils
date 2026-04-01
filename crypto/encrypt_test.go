package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	key, err := CreateKey(ctx, 32)
	require.NoError(t, err)
	plaintext := "Mary Jane Hawkins"

	ciphertext, err := Encrypt(key, []byte(plaintext))
	require.NoError(t, err)

	decryptedPlaintext, err := Decrypt(ctx, key, ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, string(decryptedPlaintext))
}
