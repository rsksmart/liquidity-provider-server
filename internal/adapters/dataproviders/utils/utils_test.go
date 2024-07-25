package utils_test

import (
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"runtime/debug"
	"testing"
	"unsafe"
)

type testStruct struct {
	AnyField string
}

func TestGetSecurePointer(t *testing.T) {
	debug.SetPanicOnFault(true)
	defer debug.SetPanicOnFault(false)
	var structToFill *testStruct

	buffer, pointer := utils.GetSecurePointer[testStruct]()
	expectedPointer := new(testStruct)
	expectedPointer.AnyField = test.AnyString
	*pointer = *expectedPointer
	structToFill = (*testStruct)(unsafe.Pointer(&buffer.Bytes()[0]))

	assert.IsType(t, expectedPointer, pointer)
	assert.Equal(t, test.AnyString, pointer.AnyField)
	assert.Equal(t, test.AnyString, structToFill.AnyField)
	assert.Equal(t, unsafe.Sizeof(expectedPointer), unsafe.Sizeof(pointer))
	assert.Equal(t, int(unsafe.Sizeof(*expectedPointer)), buffer.Size())

	_ = buffer.Seal()
	assert.Nil(t, buffer.Data())
	assert.Nil(t, buffer.Bytes())
	assert.Panics(t, func() { fmt.Println(pointer.AnyField) })
	assert.Panics(t, func() { fmt.Println(structToFill.AnyField) })
	assert.NotPanics(t, func() { fmt.Println(expectedPointer.AnyField) })
}
