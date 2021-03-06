package httpsignatures

import (
	"crypto"
	"crypto/sha512"
)

const algRsaSha512 = "RSA-SHA512"

// RsaSha512 RSA-SHA512 Algorithm
type RsaSha512 struct{}

// Algorithm Return algorithm name
func (a RsaSha512) Algorithm() string {
	return algRsaSha512
}

// Create Create signature using passed privateKey from secret
func (a RsaSha512) Create(secret Secret, data []byte) ([]byte, error) {
	return signatureRsaAlgorithmCreate(algRsaSha512, sha512.New, crypto.SHA512, secret, data)
}

// Verify Verify signature using passed publicKey from secret
func (a RsaSha512) Verify(secret Secret, data []byte, signature []byte) error {
	return signatureRsaAlgorithmVerify(algRsaSha512, sha512.New, crypto.SHA512, secret, data, signature)
}
