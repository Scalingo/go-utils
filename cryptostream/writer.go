package cryptostream

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

type EncryptedWriter struct {
	csw cipher.StreamWriter
}

func NewWriter() *EncryptedWriter {
	return &EncryptedWriter{}
}

func NewEncrypterStreamReadCloser(writer io.Writer, key []byte) (*EncryptedWriter, []byte, error) {
	iv, err := generateIV()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fail to generate IV")
	}

	if len(key) != 32 {
		return nil, nil, errors.Errorf("AES encryption key should be 32 bytes, got %d bytes", len(key))
	}

	keyBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fail to create AES cipher block")
	}

	encrypter := cipher.NewCFBEncrypter(keyBlock, iv)

	csw := cipher.StreamWriter{S: encrypter, W: writer}

	return &EncryptedWriter{csw: csw}, iv, nil
}

func (w *EncryptedWriter) Write(b []byte) (int, error) {
	return w.csw.W.Write(b)
}

func generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrapf(err, "fail to read %v bytes of random", aes.BlockSize)
	}
	return iv, nil
}
