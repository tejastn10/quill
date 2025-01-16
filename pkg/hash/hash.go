package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

// Compute a SHA-1 hash for the given data.
func ComputeSHA1(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)

	return hex.EncodeToString(hasher.Sum(nil))
}

// Compute a SHA-256 hash for the given data.
func ComputeSHA256(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)

	return hex.EncodeToString(hasher.Sum(nil))
}

// Compute a BLAKE2 hash for the given data.
func ComputeBLAKE2(data []byte) (string, error) {
	hash, err := blake2b.New256(nil)
	if err != nil {
		return "", err
	}

	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
