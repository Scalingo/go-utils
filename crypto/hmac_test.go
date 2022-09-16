package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HMAC(t *testing.T) {
	t.Parallel()

	key, err := CreateKey(128)
	require.NoError(t, err)
	plaintext := "123-12-1234"
	require.Equal(
		t,
		HMAC512(key, []byte(plaintext)),
		HMAC512(key, []byte(plaintext)),
	)
}
