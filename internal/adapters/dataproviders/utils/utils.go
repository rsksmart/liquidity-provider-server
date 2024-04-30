package utils

import (
	"github.com/awnumar/memguard"
	"unsafe"
)

func GetSecurePointer[T any]() (buffer *memguard.LockedBuffer, typePointer *T) {
	requiredTypePointer := new(T)
	lockedBuffer := memguard.NewBuffer(int(unsafe.Sizeof(*requiredTypePointer)))
	return lockedBuffer, (*T)(unsafe.Pointer(&lockedBuffer.Bytes()[0]))
}
