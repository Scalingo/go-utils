package cryptostream

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/pkg/errors"
)

type EncryptedWriter struct {
	csw             cipher.StreamWriter
	computeSHA256   bool
	plainSHA256     hash.Hash
	encryptedSHA256 hash.Hash
}

type WriterOpt func(*EncryptedWriter)

func WithSHA256() WriterOpt {
	return WriterOpt(func(w *EncryptedWriter) {
		w.computeSHA256 = true
		w.plainSHA256 = sha256.New()
		w.encryptedSHA256 = sha256.New()
	})
}

func NewWriter(writer io.Writer, key []byte) (*EncryptedWriter, []byte, error) {
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
	if w.computeSHA256 {
	}
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
