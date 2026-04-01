package crypto

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"github.com/Scalingo/go-utils/errors/v3"
)

// CreateKey creates a key of a given size by reading that much data off the crypto/rand reader.
func CreateKey(ctx context.Context, keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := cryptorand.Read(key)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "generate random bytes")
	}
	return key, nil
}

// CreateKeyString generates a new key and returns it as a hex string.
func CreateKeyString(ctx context.Context, keySize int) (string, error) {
	key, err := CreateKey(ctx, keySize)
	if err != nil {
		return "", errors.Wrap(ctx, err, "create key")
	}
	return hex.EncodeToString(key), nil
}

// CreateKeyBase64String generates a new key and returns it as a base64 std encoding string.
func CreateKeyBase64String(ctx context.Context, keySize int) (string, error) {
	key, err := CreateKey(ctx, keySize)
	if err != nil {
		return "", errors.Wrap(ctx, err, "create key")
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// ParseKey parses a key from an hexadecimal representation.
func ParseKey(ctx context.Context, key string) ([]byte, error) {
	decoded, err := hex.DecodeString(key)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "decode hexadecimal key")
	}
	if len(decoded) != DefaultKeySize {
		return nil, errors.New(ctx, "parse key: invalid key length")
	}
	return decoded, nil
}
