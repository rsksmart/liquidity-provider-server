package utils

func FirstNonZero[T comparable](values ...T) T {
	var zeroValue T
	for _, v := range values {
		if v != zeroValue {
			return v
		}
	}
	return zeroValue
}
