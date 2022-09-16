package security

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/Scalingo/go-utils/crypto"
)

// TokenGenerator lets you generate a Token.
type TokenGenerator interface {
	GenerateToken(context.Context, string) (Token, error)
}

// Token contains a hashed payload generated at a specific time.
type Token struct {
	// GeneratedAt is the token generation date represented as a Unix time
	GeneratedAt int64
	// Hash is the hex encoded HMAC of the token
	Hash string
}

type tokenGenerator struct {
	tokenSecretKey []byte
	tokenValidity  time.Duration
	now            func() time.Time
}

// NewTokenGenerator instantiates a new TokenGenerator with the given token configuration:
// - tokenSecretKeyHex: secret to generate the token.
// - tokenValidity: validity duration of the token.
func NewTokenGenerator(tokenSecretKeyHex string, tokenValidity time.Duration) (TokenGenerator, error) {
	tokenSecretKey, err := hex.DecodeString(tokenSecretKeyHex)
	if err != nil {
		return nil, errors.Wrap(err, "fail to decode the download token hex representation")
	}

	return tokenGenerator{
		tokenSecretKey: tokenSecretKey,
		tokenValidity:  tokenValidity,
		now:            time.Now,
	}, nil
}

// GenerateToken generates a new token hashed with HMAC-SHA256 for the given payload.
func (g tokenGenerator) GenerateToken(ctx context.Context, payload string) (Token, error) {
	generatedAtTimestamp := g.now().Unix()
	// Generate a hash for those metadata
	hash := crypto.HMAC256(g.tokenSecretKey, []byte(generatePlainText(generatedAtTimestamp, payload)))

	return Token{
		GeneratedAt: generatedAtTimestamp,
		Hash:        hex.EncodeToString(hash),
	}, nil
}

func generatePlainText(timestamp int64, payload string) string {
	return fmt.Sprintf("%v/%v", timestamp, payload)
}
