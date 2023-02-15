package testmocks

import (
	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
	hash        string
	quote       *pegin.Quote
	pegoutQuote *pegout.Quote
}

func NewDbMock(h string, q *pegin.Quote, pq *pegout.Quote) (*mongoDB.DB, error) {
	return nil, nil
	// return &DbMock{
	// 	hash:        h,
	// 	quote:       q,
	// 	pegoutQuote: pq,
	// }
}

func NewDbMockData(h string, q *pegin.Quote, pq *pegout.Quote) *DbMock {
	return &DbMock{
		hash:        h,
		quote:       q,
		pegoutQuote: pq,
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

func (d *DbMock) InsertQuote(id string, q *pegin.Quote) error {
	d.Called(id, q)
	return nil
}

func (d *DbMock) GetQuote(quoteHash string) (*pegin.Quote, error) {
	d.Called(quoteHash)
	return d.quote, nil
}

func (d *DbMock) RetainQuote(entry *types.RetainedQuote) error {
	d.Called(entry)
	return nil
}

func (d *DbMock) GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error) {
	d.Called(filter)
	return []*types.RetainedQuote{{QuoteHash: d.hash}}, nil
}

func (d *DbMock) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	d.Called(hash)
	return nil, nil
}

func (d *DbMock) DeleteExpiredQuotes(expTimestamp int64) error {
	d.Called(expTimestamp)
	return nil
}

func (d *DbMock) UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error {
	d.Called(hash, oldState, newState)
	return nil
}

func (d *DbMock) GetLockedLiquidity() (*types.Wei, error) {
	d.Called()
	return new(types.Wei), nil
}

func (d *DbMock) InsertPegOutQuote(id string, q *pegout.Quote, derivationAddress string) error {
	return nil
}

func (d *DbMock) GetPegOutQuote(quoteHash string) (*pegout.Quote, error) {
	d.Called(quoteHash)
	return d.pegoutQuote, nil
}

func (d *DbMock) RetainPegOutQuote(entry *pegout.RetainedQuote) error {
	return nil
}

func (d *DbMock) GetRetainedPegOutQuote(hash string) (*pegout.RetainedQuote, error) {
	d.Called(hash)
	return nil, nil
}

func (d *DbMock) UpdateRetainedPegOutQuoteState(
	hash string,
	oldState types.RQState,
	newState types.RQState) error {
	d.Called(hash, oldState, newState)
	return nil
}
