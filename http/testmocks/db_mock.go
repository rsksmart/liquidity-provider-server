package testmocks

import (
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
	hash string
	quote *types.Quote
}

func NewDbMock(h string, q *types.Quote) *DbMock {
	return &DbMock{
		hash: h,
		quote: q,
	}
}

func (d *DbMock) Close() error {
	d.Called()
	return nil
}

func (d *DbMock) InsertQuote(id string, q *types.Quote) error {
	d.Called(id, q)
	return nil
}

func (d *DbMock) GetQuote(quoteHash string) (*types.Quote, error) {
	d.Called(quoteHash)
	return d.quote, nil
}

func (d *DbMock) RetainQuote(quote *types.RetainedQuote) error {
	d.Called(quote)
	return nil
}

func (d *DbMock) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	d.Called(hash)
	return &types.RetainedQuote{ QuoteHash: hash }, nil
}

func (d *DbMock) DeleteRetainedQuote(hash string) error {
	d.Called(hash)
	return nil
}

func (d *DbMock) GetRetainedQuotes() ([]*types.RetainedQuote, error) {
	d.Called()
	return []*types.RetainedQuote{{QuoteHash: d.hash}}, nil
}

func (d *DbMock) SetRetainedQuoteCalledForUserFlag(hash string) error {
	d.Called(hash)
	return nil
}
