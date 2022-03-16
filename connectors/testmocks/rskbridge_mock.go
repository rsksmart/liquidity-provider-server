package testmocks

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/mock"
	"math/big"
)

type RskBridgeMock struct {
	mock.Mock
}

func (B *RskBridgeMock) GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	args := B.Called(opts)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (B *RskBridgeMock) GetFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	args := B.Called(opts)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (B *RskBridgeMock) GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	args := B.Called(opts)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (B *RskBridgeMock) GetFederationAddress(opts *bind.CallOpts) (string, error) {
	args := B.Called(opts)
	return args.Get(0).(string), args.Error(1)
}

func (B *RskBridgeMock) GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, arg1 string) ([]byte, error) {
	args := B.Called(opts, index, arg1)
	return args.Get(0).([]byte), args.Error(1)
}

func (B *RskBridgeMock) GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error) {
	args := B.Called(opts)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (B *RskBridgeMock) GetActivePowpegRedeemScript(opts *bind.CallOpts) ([]byte, error) {
	args := B.Called(opts)
	return args.Get(0).([]byte), args.Error(1)
}
