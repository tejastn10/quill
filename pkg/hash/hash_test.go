package hash_test

import (
	"testing"

	"github.com/tejastn10/quill/pkg/hash"
)

func TestHashFunctions(t *testing.T) {
	data := []byte("test content")
	sha256Hash := hash.ComputeSHA256(data)
	blake2Hash, err := hash.ComputeBLAKE2(data)

	if err != nil {
		t.Fatalf("error computing BLAKE2 hash: %v", err)
	}

	if sha256Hash == blake2Hash {
		t.Errorf("hash functions should produce unique results")
	}
}
