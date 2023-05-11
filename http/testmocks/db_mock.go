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

func (d *DbMock) ResetProviders(providers []*types.GlobalProvider) error {
	return d.Called(providers).Error(0)
}

func (d *DbMock) UpdateDepositedPegOutQuote(hash string, depositBlockNumber uint64) error {
	args := d.Called(hash, depositBlockNumber)
	return args.Error(0)
}

func (d *DbMock) GetRetainedPegOutQuoteByState(filter []types.RQState) ([]*pegout.RetainedQuote, error) {
	args := d.Called(filter)
	return args.Get(0).([]*pegout.RetainedQuote), args.Error(1)
}

func (d *DbMock) SaveAddressKeys(quoteHash string, addr string, pubKey []byte, privateKey []byte) error {
	args := d.Called(quoteHash, addr, pubKey, privateKey)
	return args.Error(0)
}

func (d *DbMock) GetAddressKeys(quoteHash string) (*mongoDB.PegoutKeys, error) {
	args := d.Called(quoteHash)
	return args.Get(0).(*mongoDB.PegoutKeys), args.Error(1)
}

func (d *DbMock) GetProvider(u uint64) (*mongoDB.ProviderAddress, error) {
	arg := d.Called(u)
	return arg.Get(0).(*mongoDB.ProviderAddress), arg.Error(1)
}

func (d *DbMock) GetLockedLiquidityPegOut() (uint64, error) {
	args := d.Called()
	return uint64(args.Int(0)), args.Error(1)
}

func NewDbMock(h string, q *pegin.Quote, pq *pegout.Quote) (*DbMock, error) {
	return &DbMock{
		hash:        h,
		quote:       q,
		pegoutQuote: pq,
	}, nil
}

func (d *DbMock) GetProviders() ([]*mongoDB.ProviderAddress, error) {
	args := d.Called()
	return args.Get(0).([]*mongoDB.ProviderAddress), args.Error(1)
}

func (d *DbMock) InsertProvider(id int64, address string) error {
	args := d.Called(id, address)
	return args.Error(0)
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

func (d *DbMock) InsertPegOutQuote(id string, q *pegout.Quote) error {
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
