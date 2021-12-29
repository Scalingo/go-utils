package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()

	key, err := CreateKey(32)
	require.NoError(t, err)
	plaintext := "Mary Jane Hawkins"

	ciphertext, err := Encrypt(key, []byte(plaintext))
	require.NoError(t, err)

	decryptedPlaintext, err := Decrypt(key, ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, string(decryptedPlaintext))
}
