package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"

	"github.com/Scalingo/go-utils/errors/v3"
)

// Decrypt decrypts data with the given key.
func Decrypt(ctx context.Context, key, cipherText []byte) ([]byte, error) {
	if len(cipherText) < aes.BlockSize {
		return nil, errors.Errorf(ctx, "decrypt ciphertext: input is smaller than AES block size %d", aes.BlockSize)
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)
	return cipherText, nil
}
