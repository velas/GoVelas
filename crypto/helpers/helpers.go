package helpers

import (
	"encoding/binary"
	"encoding/json"
)

// ToJSON convert object to json byte slice
func ToJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// UInt32ToBytes convert uint32 to LittleEndian byte slice
func UInt32ToBytes(val uint32) []byte {
	slice := make([]byte, 4)
	binary.LittleEndian.PutUint32(slice, val)
	return slice
}

// BytesToUInt32 convert byte slice to uint32
func BytesToUInt32(slice []byte) uint32 {
	return binary.LittleEndian.Uint32(slice)
}

// UInt64ToBytes convert uint64 to byte slice
func UInt64ToBytes(val uint64) []byte {
	slice := make([]byte, 8)
	binary.LittleEndian.PutUint64(slice, val)
	return slice
}

// BytesToUInt64 convert byte slice to uint32
func BytesToUInt64(slice []byte) uint64 {
	return binary.LittleEndian.Uint64(slice)
}
