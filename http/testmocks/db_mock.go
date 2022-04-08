package testmocks

import (
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
	hash  string
	quote *types.Quote
}

func NewDbMock(h string, q *types.Quote) *DbMock {
	return &DbMock{
		hash:  h,
		quote: q,
	}
}

func (d *DbMock) CheckConnection() error {
	args := d.Called()
	return args.Error(0)
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

func (d *DbMock) DeleteExpiredQuotes(expTimestamp int64) error {
	d.Called(expTimestamp)
	return nil
}

func (d *DbMock) RetainQuote(quote *types.RetainedQuote) error {
	d.Called(quote)
	return nil
}

func (d *DbMock) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	d.Called(hash)
	return nil, nil
}

func (d *DbMock) DeleteRetainedQuote(hash string) error {
	d.Called(hash)
	return nil
}

func (d *DbMock) GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error) {
	d.Called(filter)
	return []*types.RetainedQuote{{QuoteHash: d.hash}}, nil
}

func (d *DbMock) UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error {
	d.Called(hash, oldState, newState)
	return nil
}

func (d *DbMock) GetLockedLiquidity() (*types.Wei, error) {
	d.Called()
	return new(types.Wei), nil
}
