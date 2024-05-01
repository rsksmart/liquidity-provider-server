package utils

import (
	"github.com/awnumar/memguard"
	log "github.com/sirupsen/logrus"
	"net/http"
	"unsafe"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func CloseBodyIfExists(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}
	if err := res.Body.Close(); err != nil {
		log.Error("Error closing response body: ", err)
	}
}

func GetSecurePointer[T any]() (buffer *memguard.LockedBuffer, typePointer *T) {
	requiredTypePointer := new(T)
	lockedBuffer := memguard.NewBuffer(int(unsafe.Sizeof(*requiredTypePointer)))
	return lockedBuffer, (*T)(unsafe.Pointer(&lockedBuffer.Bytes()[0]))
}
