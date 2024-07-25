package utils

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/argon2"
)

const (
	Argon2Time                = 1
	Argon2Memory              = 64 * 1024
	Argon2Threads             = 4
	Argon2HashSize            = 32
	argon2RecommendedSaltSize = 16
)

var Argon2SaltSizeError = fmt.Errorf("the recommended size for a salt in argon2 is at least %d bytes", argon2RecommendedSaltSize)

func HashAndSaltArgon2(value string, saltSize int64) (hash []byte, salt string, err error) {
	if saltSize < argon2RecommendedSaltSize {
		return nil, "", Argon2SaltSizeError
	}
	saltBytes, err := GetRandomBytes(saltSize)
	if err != nil {
		return nil, "", err
	}
	saltValue := hex.EncodeToString(saltBytes)
	result, err := HashArgon2(value, saltValue)
	if err != nil {
		return nil, "", err
	}
	return result, hex.EncodeToString(saltBytes), nil
}

func HashArgon2(value, salt string) ([]byte, error) {
	decodedSalt, err := hex.DecodeString(salt)
	if err != nil {
		return nil, err
	} else if len(decodedSalt) < argon2RecommendedSaltSize {
		return nil, Argon2SaltSizeError
	}
	hash := argon2.IDKey([]byte(value), decodedSalt, Argon2Time, Argon2Memory, Argon2Threads, Argon2HashSize)
	return hash, nil
}
