package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

func GetRandomInt() (int64, error) {
	random, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return 0, err
	}
	return random.Int64(), nil
}

// MustGetRandomInt same as GetRandomInt but panics if error
func MustGetRandomInt() int64 {
	random, err := GetRandomInt()
	if err != nil {
		panic(err)
	}
	return random
}

func GetRandomBytes(numberOfBytes int64) ([]byte, error) {
	random := make([]byte, numberOfBytes)

	_, err := rand.Read(random)
	if err != nil {
		return nil, err
	}
	return random, nil
}

// MustGetRandomBytes same as GetRandomBytes but panics if error
func MustGetRandomBytes(numberOfBytes int64) []byte {
	random, err := GetRandomBytes(numberOfBytes)
	if err != nil {
		panic(err)
	}
	return random
}
