package cryptostream

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/pkg/errors"
)

type EncryptedReadCloser struct {
	stream *cipher.StreamReader
	closer io.Closer
}

func NewReadCloser(rc io.ReadCloser, key []byte, iv []byte) (*EncryptedReadCloser, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to create encrypted ReadCloser")
	}

	decrypter := cipher.NewCFBDecrypter(block, iv)

	return &EncryptedReadCloser{
		stream: &cipher.StreamReader{R: rc, S: decrypter},
		closer: rc,
	}, nil
}

func (reader *EncryptedReadCloser) Close() error {
	return reader.closer.Close()
}

func (reader *EncryptedReadCloser) Read(b []byte) (int, error) {
	return reader.stream.Read(b)
}
