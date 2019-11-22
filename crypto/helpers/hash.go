package helpers

var emptyHash [32]byte

// ToHash convert slice to hash
func ToHash(in []byte) [32]byte {
	var hash [32]byte
	copy(hash[:], in[:32])
	return hash
}

// HashIsEmpty check hash is empty
func HashIsEmpty(hash [32]byte) bool {
	return hash == emptyHash
}
