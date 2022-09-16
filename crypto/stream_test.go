package crypto

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Stream_EncrypterDecrypter(t *testing.T) {
	t.Parallel()

	encKey, err := CreateKey(32)
	require.NoError(t, err)
	macKey, err := CreateKey(32)
	require.NoError(t, err)
	plaintext := "Eleven is the best person in all of Hawkins Indiana. Some more text"
	pt := []byte(plaintext)

	src := bytes.NewReader(pt)

	se, err := NewStreamEncrypter(encKey, macKey, src)
	require.NoError(t, err)
	assert.NotNil(t, se)

	encrypted, err := io.ReadAll(se)
	require.NoError(t, err)
	assert.NotNil(t, encrypted)

	sd, err := NewStreamDecrypter(encKey, macKey, se.Meta(), bytes.NewReader(encrypted))
	require.NoError(t, err)
	assert.NotNil(t, sd)

	decrypted, err := io.ReadAll(sd)
	require.NoError(t, err)
	assert.Equal(t, plaintext, string(decrypted))

	require.Nil(t, sd.Authenticate())
}
