package testmocks

import (
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d DbMock) Close() error {
	d.Called()
	return nil
}

func (d DbMock) InsertQuote(id string, q *types.Quote) error {
	d.Called(id, q)
	return nil
}

func (d DbMock) GetQuote(quoteHash string) (*types.Quote, error) {
	d.Called(quoteHash)
	return &types.Quote{}, nil
}
