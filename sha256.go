package httpsignatures

import (
	"crypto/sha256"
)

const algoSha256 = "SHA-256"

// Sha256 Sha256 Algorithm
type Sha256 struct{}

// Algorithm Return algorithm name
func (a Sha256) Algorithm() string {
	return algoSha256
}

// Create Create hash
func (a Sha256) Create(data []byte) ([]byte, error) {
	return digestHashAlgorithmCreate(sha256.New, data)
}

// Verify Verify hash
func (a Sha256) Verify(data []byte, digest []byte) error {
	return digestHashAlgorithmVerify(sha256.New, data, digest)
}
