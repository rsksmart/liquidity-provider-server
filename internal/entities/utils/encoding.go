package utils

import (
	"encoding/hex"
	"fmt"
)

func DecodeKey(key string, expectedBytes int) ([]byte, error) {
	var err error
	var bytes []byte
	if bytes, err = hex.DecodeString(key); err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}
	if len(bytes) != expectedBytes {
		return nil, fmt.Errorf("key length is not %d bytes, %s is %d bytes long", expectedBytes, key, len(bytes))
	}
	return bytes, nil
}

// To32Bytes utility to convert a byte slice to a fixed size byte array, if input has
// more than 32 bytes they won't be copied.
func To32Bytes(value []byte) [32]byte {
	var bytes [32]byte
	copy(bytes[:], value)
	return bytes
}
