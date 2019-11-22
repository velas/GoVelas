// This package contains helpers for work with crypto functions
package helpers

import (
	"fmt"
)

// GetHash convert slice of bytes to limited slice
func GetHash(in []byte) ([32]byte, error) {
	var hash [32]byte

	if len(in) != 32 {
		return [32]byte{}, fmt.Errorf("getHash: invalid len of hash :%d", len(in))
	}

	copy(hash[:], in[:32])
	return hash, nil
}
