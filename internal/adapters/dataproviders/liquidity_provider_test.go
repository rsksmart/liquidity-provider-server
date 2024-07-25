package dataproviders_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"math/big"
	"strings"
	"testing"
)

const (
	rskTestAddress = "0x7c292eb881fd15605f7a85c24f4909381d36c3b9"
	quoteHash      = "5f677ed167ea3af1205ee45c64bf9883338ba9ae51f2d4e1ada949ebbff7d179"
)

func TestLocalLiquidityProvider_BtcAddress(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
	lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
	btcWallet.On("Address").Return(test.AnyAddress)
	assert.Equal(t, test.AnyAddress, lp.BtcAddress())
}

func TestLocalLiquidityProvider_RskAddress(t *testing.T) {
	signer := new(mocks.TransactionSignerMock)
	lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, signer, nil, blockchain.RskContracts{})
	signer.On("Address").Return(common.HexToAddress(rskTestAddress))
	assert.Equal(t, strings.ToLower(rskTestAddress), lp.RskAddress())
}

func TestLocalLiquidityProvider_SignQuote(t *testing.T) {
	const (
		signatureBeforeSum = "ce55f807c9f533bdf58b0bfd072dadfdd443cb521aef104f4d4014dcf4da7db418d142dfa0a26edbd169930189ed1a23b9bd8e09c7b01f3832e26fc7855f89a900"
		signatureAfterSum  = "ce55f807c9f533bdf58b0bfd072dadfdd443cb521aef104f4d4014dcf4da7db418d142dfa0a26edbd169930189ed1a23b9bd8e09c7b01f3832e26fc7855f89a91b"
	)
	var buffer bytes.Buffer
	hashBytes, err := hex.DecodeString(quoteHash)
	require.NoError(t, err)
	buffer.WriteString("\x19Ethereum Signed Message:\n32")
	buffer.Write(hashBytes)
	signer := new(mocks.TransactionSignerMock)
	signatureBytes, err := hex.DecodeString(signatureBeforeSum)
	require.NoError(t, err)
	signer.On("SignBytes", mock.MatchedBy(func(content []byte) bool {
		return bytes.Equal(content, crypto.Keccak256(buffer.Bytes()))
	})).Return(signatureBytes, nil)
	lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, signer, nil, blockchain.RskContracts{})
	result, err := lp.SignQuote(quoteHash)
	signer.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, signatureAfterSum, result)
}

func TestLocalLiquidityProvider_SignQuote_ErrorHandling(t *testing.T) {
	t.Run("Invalid hash", func(t *testing.T) {
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		result, err := lp.SignQuote(test.AnyString)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Signing error", func(t *testing.T) {
		signer := new(mocks.TransactionSignerMock)
		signer.On("SignBytes", mock.Anything).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, signer, nil, blockchain.RskContracts{})
		result, err := lp.SignQuote(quoteHash)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestLocalLiquidityProvider_AvailablePegoutLiquidity(t *testing.T) {
	pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateWaitingForDepositConfirmations,
	).Return([]quote.RetainedPegoutQuote{
		{RequiredLiquidity: entities.NewWei(100)},
		{RequiredLiquidity: entities.NewWei(300)},
		{RequiredLiquidity: entities.NewWei(200)},
	}, nil).Once()
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(1500), nil).Once()
	lp := dataproviders.NewLocalLiquidityProvider(nil, pegoutRepository, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
	liquidity, err := lp.AvailablePegoutLiquidity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(900), liquidity)
	btcWallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

func TestLocalLiquidityProvider_AvailablePegoutLiquidity_ErrorHandling(t *testing.T) {
	t.Run("Error getting btc wallet balance when checking available pegout liquidity", func(t *testing.T) {
		btcWallet := new(mocks.BitcoinWalletMock)
		btcWallet.On("GetBalance").Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
		liquidity, err := lp.AvailablePegoutLiquidity(context.Background())
		require.Error(t, err)
		assert.Nil(t, liquidity)
	})
	t.Run("Error getting pegout quotes from db when checking available pegout liquidity", func(t *testing.T) {
		pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
		pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx,
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
		).Return(nil, assert.AnError).Once()
		btcWallet := new(mocks.BitcoinWalletMock)
		btcWallet.On("GetBalance").Return(entities.NewWei(500), nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, pegoutRepository, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
		liquidity, err := lp.AvailablePegoutLiquidity(context.Background())
		require.Error(t, err)
		assert.Nil(t, liquidity)
	})
}

func TestLocalLiquidityProvider_HasPegoutLiquidity(t *testing.T) {
	pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateWaitingForDepositConfirmations,
	).Return([]quote.RetainedPegoutQuote{
		{RequiredLiquidity: entities.NewWei(100)},
		{RequiredLiquidity: entities.NewWei(200)},
		{RequiredLiquidity: entities.NewWei(150)},
	}, nil).Times(3)
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(500), nil).Times(3)
	lp := dataproviders.NewLocalLiquidityProvider(nil, pegoutRepository, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
	testCases := []struct {
		amount        *entities.Wei
		expectedError string
	}{
		{amount: entities.NewWei(50), expectedError: ""},
		{amount: entities.NewWei(150), expectedError: "not enough liquidity"},
		{amount: entities.NewWei(20), expectedError: ""},
	}
	for _, tc := range testCases {
		err := lp.HasPegoutLiquidity(context.Background(), tc.amount)
		if tc.expectedError == "" {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		}
	}
	btcWallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

func TestLocalLiquidityProvider_HasPegoutLiquidity_ErrorHandling(t *testing.T) {
	t.Run("Error getting btc wallet balance", func(t *testing.T) {
		btcWallet := new(mocks.BitcoinWalletMock)
		btcWallet.On("GetBalance").Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
		err := lp.HasPegoutLiquidity(context.Background(), entities.NewWei(1))
		require.Error(t, err)
	})
	t.Run("Error getting pegout quotes from db", func(t *testing.T) {
		pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
		pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx,
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
		).Return(nil, assert.AnError).Times(3)
		btcWallet := new(mocks.BitcoinWalletMock)
		btcWallet.On("GetBalance").Return(entities.NewWei(500), nil).Times(3)
		lp := dataproviders.NewLocalLiquidityProvider(nil, pegoutRepository, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{})
		err := lp.HasPegoutLiquidity(context.Background(), entities.NewWei(1))
		require.Error(t, err)
	})
}

func TestLocalLiquidityProvider_HasPeginLiquidity(t *testing.T) {
	signer := new(mocks.TransactionSignerMock)
	signer.On("Address").Return(common.HexToAddress(rskTestAddress)).Times(6)
	peginRepository := new(mocks.PeginQuoteRepositoryMock)
	peginRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PeginStateWaitingForDeposit,
	).Return([]quote.RetainedPeginQuote{
		{RequiredLiquidity: entities.NewWei(100)},
		{RequiredLiquidity: entities.NewWei(200)},
		{RequiredLiquidity: entities.NewWei(150)},
	}, nil).Times(3)
	pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PegoutStateRefundPegOutSucceeded,
	).Return([]quote.RetainedPegoutQuote{
		{RequiredLiquidity: entities.NewWei(30)},
		{RequiredLiquidity: entities.NewWei(50)},
	}, nil).Times(3)
	lbcMock := new(mocks.LbcMock)
	lbcMock.On("GetBalance", rskTestAddress).Return(entities.NewWei(400), nil).Times(3)
	rpcMock := new(mocks.RootstockRpcServerMock)
	rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(300), nil).Times(3)
	lp := dataproviders.NewLocalLiquidityProvider(peginRepository, pegoutRepository, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})
	testCases := []struct {
		amount        *entities.Wei
		expectedError string
	}{
		{amount: entities.NewWei(170), expectedError: ""},
		{amount: entities.NewWei(200), expectedError: "not enough liquidity"},
		{amount: entities.NewWei(50), expectedError: ""},
	}
	for _, tc := range testCases {
		err := lp.HasPeginLiquidity(context.Background(), tc.amount)
		if tc.expectedError == "" {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		}
	}
	lbcMock.AssertExpectations(t)
	rpcMock.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	signer.AssertExpectations(t)
}

func TestLocalLiquidityProvider_HasPeginLiquidity_ErrorHandling(t *testing.T) {
	signer := new(mocks.TransactionSignerMock)
	signer.On("Address").Return(common.HexToAddress(rskTestAddress))
	t.Run("Error getting balance from RSK RPC server", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{})
		err := lp.HasPeginLiquidity(context.Background(), entities.NewWei(1))
		require.Error(t, err)
	})
	t.Run("Error getting balance from LBC", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(100), nil).Once()
		lbcMock := new(mocks.LbcMock)
		lbcMock.On("GetBalance", rskTestAddress).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})
		err := lp.HasPeginLiquidity(context.Background(), entities.NewWei(1))
		require.Error(t, err)
	})
	t.Run("Error pegin quotes from db", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(100), nil).Once()
		lbcMock := new(mocks.LbcMock)
		lbcMock.On("GetBalance", rskTestAddress).Return(entities.NewWei(200), nil).Once()
		peginRepository := new(mocks.PeginQuoteRepositoryMock)
		peginRepository.On("GetRetainedQuoteByState", test.AnyCtx,
			quote.PeginStateWaitingForDeposit,
		).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(peginRepository, nil, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})
		err := lp.HasPeginLiquidity(context.Background(), entities.NewWei(1))
		require.Error(t, err)
	})
}

func TestLocalLiquidityProvider_AvailablePeginLiquidity(t *testing.T) {
	signer := new(mocks.TransactionSignerMock)
	signer.On("Address").Return(common.HexToAddress(rskTestAddress)).Twice()
	peginRepository := new(mocks.PeginQuoteRepositoryMock)
	peginRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PeginStateWaitingForDeposit,
	).Return([]quote.RetainedPeginQuote{
		{RequiredLiquidity: entities.NewWei(300)},
		{RequiredLiquidity: entities.NewWei(500)},
		{RequiredLiquidity: entities.NewWei(400)},
	}, nil).Once()
	pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PegoutStateRefundPegOutSucceeded,
	).Return([]quote.RetainedPegoutQuote{
		{RequiredLiquidity: entities.NewWei(100)},
		{RequiredLiquidity: entities.NewWei(150)},
	}, nil).Once()
	lbcMock := new(mocks.LbcMock)
	lbcMock.On("GetBalance", rskTestAddress).Return(entities.NewWei(2000), nil).Once()
	rpcMock := new(mocks.RootstockRpcServerMock)
	rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(800), nil).Once()
	lp := dataproviders.NewLocalLiquidityProvider(peginRepository, pegoutRepository, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})

	liquidity, err := lp.AvailablePeginLiquidity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(1350), liquidity)

	lbcMock.AssertExpectations(t)
	rpcMock.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	signer.AssertExpectations(t)
}

func TestLocalLiquidityProvider_AvailablePeginLiquidity_ErrorHandling(t *testing.T) {
	signer := new(mocks.TransactionSignerMock)
	signer.On("Address").Return(common.HexToAddress(rskTestAddress))
	t.Run("Error getting balance from RSK RPC server when getting available pegin liquidity", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{})
		liquidity, err := lp.AvailablePeginLiquidity(context.Background())
		require.Error(t, err)
		assert.Nil(t, liquidity)
	})
	t.Run("Error getting balance from LBC when getting available pegin liquidity", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(100), nil).Once()
		lbcMock := new(mocks.LbcMock)
		lbcMock.On("GetBalance", rskTestAddress).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})
		liquidity, err := lp.AvailablePeginLiquidity(context.Background())
		require.Error(t, err)
		assert.Nil(t, liquidity)
	})
	t.Run("Error getting pegin quotes from db when getting available pegin liquidity", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(100), nil).Once()
		lbcMock := new(mocks.LbcMock)
		lbcMock.On("GetBalance", rskTestAddress).Return(entities.NewWei(200), nil).Once()
		peginRepository := new(mocks.PeginQuoteRepositoryMock)
		peginRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateWaitingForDeposit).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(peginRepository, nil, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})
		liquidity, err := lp.AvailablePeginLiquidity(context.Background())
		require.Error(t, err)
		assert.Nil(t, liquidity)
	})
	t.Run("Error getting pegout quotes from db when getting available pegin liquidity", func(t *testing.T) {
		rpcMock := new(mocks.RootstockRpcServerMock)
		rpcMock.On("GetBalance", test.AnyCtx, rskTestAddress).Return(entities.NewWei(100), nil).Once()
		lbcMock := new(mocks.LbcMock)
		lbcMock.On("GetBalance", rskTestAddress).Return(entities.NewWei(200), nil).Once()
		peginRepository := new(mocks.PeginQuoteRepositoryMock)
		peginRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateWaitingForDeposit).
			Return([]quote.RetainedPeginQuote{{RequiredLiquidity: entities.NewWei(300)}}, nil).Once()
		pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
		pegoutRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PegoutStateRefundPegOutSucceeded).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(peginRepository, pegoutRepository, nil, blockchain.Rpc{Rsk: rpcMock}, signer, nil, blockchain.RskContracts{Lbc: lbcMock})
		liquidity, err := lp.AvailablePeginLiquidity(context.Background())
		require.Error(t, err)
		assert.Nil(t, liquidity)
	})
}

func TestLocalLiquidityProvider_GeneralConfiguration(t *testing.T) {
	message := make([]byte, 1024)
	account := test.OpenWalletForTest(t, "general-configuration")
	wallet := rootstock.NewRskWalletImpl(&rootstock.RskClient{}, account, 31)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	t.Run("Return signed general configuration from db", func(t *testing.T) {
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(getGeneralConfigurationMock(), nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.GeneralConfiguration(context.Background())
		assert.Equal(t, getGeneralConfigurationMock().Value, result)
		_, err := buff.Read(message)
		require.ErrorIs(t, err, io.EOF)
		assert.NotEqual(t, liquidity_provider.DefaultGeneralConfiguration(), result)
	})
	t.Run("Return default general configuration on db read error", func(t *testing.T) {
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		config := lp.GeneralConfiguration(context.Background())
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Error getting general configuration")
		assert.Equal(t, liquidity_provider.DefaultGeneralConfiguration(), config)
	})
	t.Run("Return default general configuration when db doesn't have configuration", func(t *testing.T) {
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(nil, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		config := lp.GeneralConfiguration(context.Background())
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Custom general configuration not found")
		assert.Equal(t, liquidity_provider.DefaultGeneralConfiguration(), config)
	})
	t.Run("Return default general configuration when db configuration is tampered", func(t *testing.T) {
		tamperedConfig := getGeneralConfigurationMock()
		tamperedConfig.Value.RskConfirmations[2000000000000000000] = 40
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(tamperedConfig, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.GeneralConfiguration(context.Background())
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Tampered general configuration")
		assert.Equal(t, liquidity_provider.DefaultGeneralConfiguration(), result)
	})
	t.Run("Return default general configuration when db configuration doesn't have valid signature", func(t *testing.T) {
		invalidSignatureConfig := getGeneralConfigurationMock()
		invalidSignatureConfig.Signature = "94530cf2d078ce7e44b4ce1d63a0cf7a225f07d4414f4dcf132f097fd027c08c7252b012ffff6855400fbc96939662904b22ce0b7a010bcb0b7a2c7db9dc26b702"
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(invalidSignatureConfig, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.GeneralConfiguration(context.Background())
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Invalid general configuration signature")
		assert.Equal(t, liquidity_provider.DefaultGeneralConfiguration(), result)
	})
}

func TestLocalLiquidityProvider_PeginConfiguration(t *testing.T) {
	message := make([]byte, 1024)
	account := test.OpenWalletForTest(t, "pegin-configuration")
	wallet := rootstock.NewRskWalletImpl(&rootstock.RskClient{}, account, 31)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	t.Run("Return signed pegin configuration from db", func(t *testing.T) {
		lpRepository.On("GetPeginConfiguration", test.AnyCtx).Return(getPeginConfigurationMock(), nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.PeginConfiguration(context.Background())
		assert.Equal(t, getPeginConfigurationMock().Value, result)
		_, err := buff.Read(message)
		require.ErrorIs(t, err, io.EOF)
		assert.NotEqual(t, liquidity_provider.DefaultPeginConfiguration(), result)
	})
	t.Run("Return default pegin configuration on db read error", func(t *testing.T) {
		lpRepository.On("GetPeginConfiguration", test.AnyCtx).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		config := lp.PeginConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPeginConfiguration(), config)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Error getting pegin configuration")
	})
	t.Run("Return default pegin configuration when db doesn't have configuration", func(t *testing.T) {
		lpRepository.On("GetPeginConfiguration", test.AnyCtx).Return(nil, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		config := lp.PeginConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPeginConfiguration(), config)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Custom pegin configuration not found")
	})
	t.Run("Return default pegin configuration when db configuration is tampered", func(t *testing.T) {
		tamperedConfig := getPeginConfigurationMock()
		tamperedConfig.Value.MinValue = entities.NewWei(1)
		lpRepository.On("GetPeginConfiguration", test.AnyCtx).Return(tamperedConfig, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.PeginConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPeginConfiguration(), result)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Tampered pegin configuration")
	})
	t.Run("Return default pegin configuration when db configuration doesn't have valid signature", func(t *testing.T) {
		invalidSignatureConfig := getPeginConfigurationMock()
		invalidSignatureConfig.Signature = "93530cf2d078ce7e44c4ce1d63a0cf7a225f07d4414f4dcf132f097fd027c08c7252b012f1ff6855400fbc96939662904b22ce0b7a010bcb0b7a2c7db9dc26b702"
		lpRepository.On("GetPeginConfiguration", test.AnyCtx).Return(invalidSignatureConfig, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.PeginConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPeginConfiguration(), result)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Invalid pegin configuration signature")
	})
}

func TestLocalLiquidityProvider_PegoutConfiguration(t *testing.T) {
	message := make([]byte, 1024)
	account := test.OpenWalletForTest(t, "pegout-configuration")
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	wallet := rootstock.NewRskWalletImpl(&rootstock.RskClient{}, account, 31)
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	t.Run("Return signed pegout configuration from db", func(t *testing.T) {
		lpRepository.On("GetPegoutConfiguration", test.AnyCtx).Return(getPegoutConfigurationMock(), nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.PegoutConfiguration(context.Background())
		assert.Equal(t, getPegoutConfigurationMock().Value, result)
		_, err := buff.Read(message)
		require.ErrorIs(t, err, io.EOF)
		assert.NotEqual(t, liquidity_provider.DefaultPegoutConfiguration(), result)
	})
	t.Run("Return default pegout configuration on db read error", func(t *testing.T) {
		lpRepository.On("GetPegoutConfiguration", test.AnyCtx).Return(nil, assert.AnError).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		config := lp.PegoutConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPegoutConfiguration(), config)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Error getting pegout configuration")
	})
	t.Run("Return default pegout configuration when db doesn't have configuration", func(t *testing.T) {
		lpRepository.On("GetPegoutConfiguration", test.AnyCtx).Return(nil, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, nil, nil, blockchain.RskContracts{})
		config := lp.PegoutConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPegoutConfiguration(), config)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Custom pegout configuration not found")
	})
	t.Run("Return default pegout configuration when db configuration is tampered", func(t *testing.T) {
		tamperedConfig := getPegoutConfigurationMock()
		tamperedConfig.Value.MaxValue = entities.NewWei(1)
		lpRepository.On("GetPegoutConfiguration", test.AnyCtx).Return(tamperedConfig, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.PegoutConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPegoutConfiguration(), result)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Tampered pegout configuration")
	})
	t.Run("Return default pegout configuration when db configuration doesn't have valid signature", func(t *testing.T) {
		invalidSignatureConfig := getPegoutConfigurationMock()
		invalidSignatureConfig.Signature = "93530cf2d078ce7e44c4ce1d63a0cf7a225f07d4414f4dcf133f097fd027d08c7252b012f1ff6855400fbc96939662904b22ce0b7a010bcb0b7a2c7db9dc26b702"
		lpRepository.On("GetPegoutConfiguration", test.AnyCtx).Return(invalidSignatureConfig, nil).Once()
		lp := dataproviders.NewLocalLiquidityProvider(nil, nil, lpRepository, blockchain.Rpc{}, wallet, nil, blockchain.RskContracts{})
		result := lp.PegoutConfiguration(context.Background())
		assert.Equal(t, liquidity_provider.DefaultPegoutConfiguration(), result)
		_, err := buff.Read(message)
		buff.Reset()
		require.NoError(t, err)
		assert.Contains(t, string(message), "Invalid pegout configuration signature")
	})
}

func getGeneralConfigurationMock() *entities.Signed[liquidity_provider.GeneralConfiguration] {
	return &entities.Signed[liquidity_provider.GeneralConfiguration]{
		Value: liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				2000000000000000000: 15,
				400000000000000000:  10,
				4000000000000000000: 20,
				8000000000000000000: 25,
				9000000000000000000: 30,
				100000000000000000:  5,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				100000000000000000:  2,
				2000000000000000000: 10,
				400000000000000000:  6,
				4000000000000000000: 20,
				8000000000000000000: 40,
				9000000000000000001: 45,
			},
			PublicLiquidityCheck: false,
		},
		Signature: "12f9530beed2220769a3867a01ad7164af2d159cc93644dc8097a736f136b4ac076227ca370de81d0d66b962ac4c6f6f13920afec2919c1f9ee17c954a8690e601",
		Hash:      "83cb825a5f8dcf1bdd3cd33effffda7a34ed8b0d80a39445049ddc9c06ecb1a8",
	}
}
func getPeginConfigurationMock() *entities.Signed[liquidity_provider.PeginConfiguration] {
	maxBigInt := new(big.Int)
	maxBigInt.SetString("10000000000000000000", 10)
	return &entities.Signed[liquidity_provider.PeginConfiguration]{
		Value: liquidity_provider.PeginConfiguration{
			TimeForDeposit: 3600,
			CallTime:       7212,
			PenaltyFee:     entities.NewWei(1000000000000000),
			CallFee:        entities.NewWei(10000000000000000),
			MaxValue:       entities.NewBigWei(maxBigInt),
			MinValue:       entities.NewWei(600000000000000000),
		},
		Signature: "671242c9fe5e1183d68e7c45face4ab1be9f304ac661aa778b81fceca000d19360eadf8a9e8ddd1528f485e64cd27d18de989f3d61c55e7ca074df2312623f2f01",
		Hash:      "609a64fa0aa6a60dd2eee3bf704d4afdc00b69768e271bec72d0c576ff56baf2",
	}
}
func getPegoutConfigurationMock() *entities.Signed[liquidity_provider.PegoutConfiguration] {
	maxBigInt := new(big.Int)
	maxBigInt.SetString("10000000000000000000", 10)
	return &entities.Signed[liquidity_provider.PegoutConfiguration]{
		Value: liquidity_provider.PegoutConfiguration{
			TimeForDeposit:       3655,
			ExpireTime:           7201,
			PenaltyFee:           entities.NewWei(1000000000000000),
			CallFee:              entities.NewWei(10000000000000000),
			MaxValue:             entities.NewBigWei(maxBigInt),
			MinValue:             entities.NewWei(600000000000000000),
			ExpireBlocks:         500,
			BridgeTransactionMin: entities.NewWei(1500000000000000000),
		},
		Signature: "e3cd6e7dbf9708158ad03fb51ad493310bb93c4a5d56f6192531bb746f42251301f178ffd2173552181a729c2bb23ab98c9f3786bb7075f416bfe528f3a77d1601",
		Hash:      "b79d03455e8407660c123700e24af0ea92f6ca83ccb8918e8c155295fc9ac48d",
	}
}
