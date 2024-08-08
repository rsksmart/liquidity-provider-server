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

func GetRandomBytes(numberOfBytes int64) ([]byte, error) {
	random := make([]byte, numberOfBytes)

	_, err := rand.Read(random)
	if err != nil {
		return nil, err
	}
	return random, nil
}
