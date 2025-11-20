package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewResignationHandler(t *testing.T) {
	const (
		method = "POST"
		path   = "/providers/resignation"
	)
	t.Run("should return 204 on successful resignation", func(t *testing.T) {
		useCase := new(mocks.ResignUseCaseMock)
		useCase.EXPECT().Run().Return(nil).Once()
		handler := NewResignationHandler(useCase)
		assert.HTTPStatusCode(t, handler, method, path, nil, http.StatusNoContent)
	})
	t.Run("should return 503 when protocol is paused", func(t *testing.T) {
		useCase := new(mocks.ResignUseCaseMock)
		useCase.EXPECT().Run().Return(blockchain.ContractPausedError).Once()
		handler := NewResignationHandler(useCase)
		assert.HTTPStatusCode(t, handler, method, path, nil, http.StatusServiceUnavailable)
	})
	t.Run("should return 500 on unknown error", func(t *testing.T) {
		useCase := new(mocks.ResignUseCaseMock)
		useCase.EXPECT().Run().Return(assert.AnError).Once()
		handler := NewResignationHandler(useCase)
		assert.HTTPStatusCode(t, handler, method, path, nil, http.StatusInternalServerError)
	})
}
