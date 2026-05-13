package idempotency

type Generator interface {
    Generate(payload []byte) string
}