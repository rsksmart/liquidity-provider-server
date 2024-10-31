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
		return nil, fmt.Errorf("key length is not %d bytes, key is %d bytes long", expectedBytes, len(bytes))
	}
	return bytes, nil
}
