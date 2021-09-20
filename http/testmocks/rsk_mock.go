package testmocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/mock"
)

type RskMock struct {
	mock.Mock
}

func (m *RskMock) GetAvailableLiquidity(addr string) (*big.Int, error) {
	m.Called(addr)
	return nil, nil
}

func (m *RskMock) GetCollateral(addr string) (*big.Int, *big.Int, error) {
	m.Called(addr)
	return big.NewInt(10), big.NewInt(10), nil
}

func (m *RskMock) RegisterProvider(opts *bind.TransactOpts) error {
	m.Called()
	return nil
}

func (m *RskMock) AddCollateral(opts *bind.TransactOpts) error {
	m.Called()
	return nil
}

func (m *RskMock) GetRequiredBridgeConfirmations() int64 {
	m.Called()
	return 0
}

func (m *RskMock) GetChainId() *big.Int {
	m.Called()
	return big.NewInt(0)
}

func (m *RskMock) ParseQuote(q *types.Quote) (bindings.LiquidityBridgeContractQuote, error) {
	m.Called(q)
	return bindings.LiquidityBridgeContractQuote{}, nil
}

func (m *RskMock) RegisterPegIn(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote, signature []byte, btcRawTrx []byte, partialMerkleTree []byte, height *big.Int) (*gethTypes.Transaction, error) {
	m.Called(opt, q, signature, btcRawTrx, partialMerkleTree, height)
	return nil, nil
}

func (m *RskMock) CallForUser(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote) (*gethTypes.Transaction, error) {
	m.Called(opt, q)
	return nil, nil
}

func (m *RskMock) Connect(endpoint string) error {
	m.Called(endpoint)
	return nil
}
func (m *RskMock) Close() {
	m.Called()
}

func (m *RskMock) EstimateGas(addr string, value big.Int, data []byte) (uint64, error) {
	args := m.Called(addr, value, data)
	return args.Get(0).(uint64), nil
}

func (m *RskMock) GasPrice() (*big.Int, error) {
	args := m.Called()
	return args.Get(0).(*big.Int), nil
}
func (m *RskMock) HashQuote(q *types.Quote) (string, error) {
	args := m.Called(q)
	return args.String(), nil
}
func (m *RskMock) GetFedSize() (int, error) {
	args := m.Called()
	return args.Int(0), nil
}
func (m *RskMock) GetFedThreshold() (int, error) {
	args := m.Called()
	return args.Int(0), nil
}

func (m *RskMock) GetFedPublicKey(index int) (string, error) {
	args := m.Called(index)
	return args.String(), nil
}
func (m *RskMock) GetFedAddress() (string, error) {
	args := m.Called()
	return args.String(), nil
}
func (m *RskMock) GetActiveFederationCreationBlockHeight() (int, error) {
	args := m.Called()
	return args.Int(0), nil
}

func (m *RskMock) GetLBCAddress() string {
	args := m.Called()
	return args.String()
}
