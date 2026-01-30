package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRequestID() (requestID string) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		bytes = []byte("fallback-request-id-error")
	}
	requestID = hex.EncodeToString(bytes)
	return
}
