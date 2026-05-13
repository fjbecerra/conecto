package idempotency

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashGenerator struct{}

func (g *HashGenerator) Generate(payload []byte) string {
    h := sha256.New()
    h.Write(payload)
    return hex.EncodeToString(h.Sum(nil))
}