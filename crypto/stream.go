package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/Scalingo/go-utils/errors/v3"
)

// NewStreamEncrypter creates a new stream encrypter
func NewStreamEncrypter(ctx context.Context, encKey, macKey []byte, plainText io.Reader) (*StreamEncrypter, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create AES cipher")
	}
	iv := make([]byte, block.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "generate random IV")
	}
	stream := cipher.NewCTR(block, iv)
	mac := hmac.New(sha256.New, macKey)
	return &StreamEncrypter{
		Source: plainText,
		Block:  block,
		Stream: stream,
		Mac:    mac,
		IV:     iv,
	}, nil
}

// NewStreamDecrypter creates a new stream decrypter
func NewStreamDecrypter(ctx context.Context, encKey, macKey []byte, meta StreamMeta, cipherText io.Reader) (*StreamDecrypter, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create AES cipher")
	}
	stream := cipher.NewCTR(block, meta.IV)
	mac := hmac.New(sha256.New, macKey)
	return &StreamDecrypter{
		Source: cipherText,
		Block:  block,
		Stream: stream,
		Mac:    mac,
		Meta:   meta,
	}, nil
}

// StreamEncrypter is an encrypter for a stream of data with authentication
type StreamEncrypter struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
	Mac    hash.Hash
	IV     []byte
}

// StreamDecrypter is a decrypter for a stream of data with authentication
type StreamDecrypter struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
	Mac    hash.Hash
	Meta   StreamMeta
}

// Read encrypts the bytes of the inner reader and places them into p
func (s *StreamEncrypter) Read(p []byte) (int, error) {
	ctx := context.Background()

	n, readErr := s.Source.Read(p)
	if n > 0 {
		s.Stream.XORKeyStream(p[:n], p[:n])
		err := writeHash(ctx, s.Mac, p[:n])
		if err != nil {
			return n, errors.Wrap(ctx, err, "write stream hash")
		}
		return n, readErr
	}
	return 0, io.EOF
}

// Meta returns the encrypted stream metadata for use in decrypting. This should only be called after the stream is finished
func (s *StreamEncrypter) Meta() StreamMeta {
	return StreamMeta{IV: s.IV, Hash: s.Mac.Sum(nil)}
}

// Read reads bytes from the underlying reader and then decrypts them
func (s *StreamDecrypter) Read(p []byte) (int, error) {
	ctx := context.Background()

	n, readErr := s.Source.Read(p)
	if n > 0 {
		err := writeHash(ctx, s.Mac, p[:n])
		if err != nil {
			return n, errors.Wrap(ctx, err, "write stream hash")
		}
		s.Stream.XORKeyStream(p[:n], p[:n])
		return n, readErr
	}
	return 0, io.EOF
}

// Authenticate verifys that the hash of the stream is correct. This should only be called after processing is finished
func (s *StreamDecrypter) Authenticate(ctx context.Context) error {
	if !hmac.Equal(s.Meta.Hash, s.Mac.Sum(nil)) {
		return errors.New(ctx, "authenticate stream: hash mismatch")
	}
	return nil
}

func writeHash(ctx context.Context, mac hash.Hash, p []byte) error {
	m, err := mac.Write(p)
	if err != nil {
		return errors.Wrap(ctx, err, "write MAC data")
	}
	if m != len(p) {
		return errors.New(ctx, "write MAC data: short write")
	}
	return nil
}

// StreamMeta is metadata about an encrypted stream
type StreamMeta struct {
	// IV is the initial value for the crypto function
	IV []byte
	// Hash is the sha256 hmac of the stream
	Hash []byte
}
