package state

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidState = errors.New("invalid oauth state")
	ErrExpiredState = errors.New("expired oauth state")
)


type HMACStateSigner struct {
	Secret []byte
	TTL    time.Duration
}

func NewHMACStateSigner(secret []byte,ttl time.Duration) *HMACStateSigner {
	return &HMACStateSigner{
		Secret: secret,
		TTL: ttl,
	}
}

func (s *HMACStateSigner) Sign(connectionID string) (string, error) {


	nonce, err := randomString(16)

	if err != nil {
		return "", err
	}


	payload := StatePayload{

		ConnectionID: connectionID,

		ExpiresAt:
			time.Now().Add(s.TTL),

		Nonce: nonce,
	}


	raw, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}


	encodedPayload := base64.RawURLEncoding.EncodeToString(raw)


	signature := s.sign(encodedPayload)


	return encodedPayload + "." + signature, nil
}

func (s *HMACStateSigner) Verify(state string) (string, error) {


	parts := strings.Split(state, ".")

	if len(parts) != 2 {
		return "", ErrInvalidState
	}


	payloadEncoded := parts[0]

	signature := parts[1]


	expected :=
		s.sign(payloadEncoded)


	if !hmac.Equal(
		[]byte(signature),
		[]byte(expected),
	) {
		return "", ErrInvalidState
	}


	raw, err :=
		base64.RawURLEncoding.DecodeString(
			payloadEncoded,
		)

	if err != nil {
		return "", ErrInvalidState
	}


	var payload StatePayload


	err = json.Unmarshal(raw,&payload)

	if err != nil {
		return "", ErrInvalidState
	}


	if time.Now().After(payload.ExpiresAt) {

		return "", ErrExpiredState
	}


	return payload.ConnectionID, nil
}

func (s *HMACStateSigner) sign(payload string) string {


	mac := hmac.New(
			sha256.New,
			s.Secret,
		)


	mac.Write([]byte(payload))


	return base64.RawURLEncoding.EncodeToString(
		mac.Sum(nil),
	)
}

func randomString(size int) (string, error) {

	b := make([]byte, size)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}


	return base64.RawURLEncoding.EncodeToString(b), nil
}