package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var AnyCtx = mock.AnythingOfType("context.backgroundCtx")

const (
	AnyAddress = "any address"
	AnyString  = "any value"
)

type Case[V, R any] struct {
	Value  V
	Result R
}

type Table[V, R any] []Case[V, R]

func RunTable[V, R any](t *testing.T, table Table[V, R], validationFunction func(V) R) {
	var result R
	for _, testCase := range table {
		result = validationFunction(testCase.Value)
		assert.Equal(t, testCase.Result, result)
	}
}
