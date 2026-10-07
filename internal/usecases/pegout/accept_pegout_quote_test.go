package pegout_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var acceptPegoutQuoteHash = "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf"
var acceptPegoutQuoteHashSignature = "b062b09f5f3000f1092e606e90fa449e8527fb1bac20ff72897fd1d0a8aa3b18049d39d1956110992de0284e6d85223d4f69ed06e57184cad13abca7b421d6e41b"
var acceptPegoutQuoteEip712Hash = "95c3ca51e1abd141bed5fb1c1802236aef4e5982ffc073d4e7eed4c73d553f9a"
var ownerAccountAddress = "0x57f9F71E683E2A8ff3d2f394aE45C58b2d913A35"

func acceptTestQuote(now time.Time) quote.PegoutQuote {
	return quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(5),
		PenaltyFee:            entities.NewWei(1),
		Nonce:                 1,
		DepositAddress:        "address",
		Value:                 entities.NewWei(12),
		AgreementTimestamp:    uint32(now.Unix()),
		DepositDateLimit:      uint32(now.Unix() + 600),
		DepositConfirmations:  1,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            uint32(now.Unix() + 600),
		ExpireBlock:           1,
		GasFee:                entities.NewWei(6),
		ChainId:               31,
	}
}

func TestAcceptQuoteUseCase_Run_Paused(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	lp := new(mocks.ProviderMock)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	pegoutContract.EXPECT().GetAddress().Return("test-contract")

	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, contracts, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")
	assert.Empty(t, result)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

func TestAcceptQuoteUseCase_Run(t *testing.T) {
	quoteHash := acceptPegoutQuoteHash
	quoteMock := acceptTestQuote(time.Now())
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("GetAddress").Return("0xabcd01").Once()
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("SignPegoutQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(acceptPegoutQuoteHashSignature, nil)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), quoteHash, "")

	quoteRepositoryMock.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepositoryMock.AssertNotCalled(t, "InsertRetainedQuote")
	lp.AssertNotCalled(t, "HasPegoutLiquidity")
	require.NoError(t, err)
	assert.Equal(t, "0xabcd01", result.DepositAddress)
	assert.Equal(t, acceptPegoutQuoteHashSignature, result.Signature)
}

func TestAcceptQuoteUseCase_Run_AlreadyAccepted(t *testing.T) {
	quoteMock := acceptTestQuote(time.Now())
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         acceptPegoutQuoteHash,
		DepositAddress:    "0xexisting",
		Signature:         "existing-sig",
		RequiredLiquidity: entities.NewWei(18),
		State:             quote.PegoutStateClaimed,
	}
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&quoteMock, nil).Once()
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&retainedQuote, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	require.NoError(t, err)
	assert.Equal(t, retainedQuote.Signature, result.Signature)
	assert.Equal(t, retainedQuote.DepositAddress, result.DepositAddress)
	lp.AssertNotCalled(t, "SignPegoutQuote")
	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
}

func TestAcceptQuoteUseCase_Run_QuoteNotFound(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
}

func TestAcceptQuoteUseCase_Run_ExpiredQuote(t *testing.T) {
	now := time.Now()
	expired := acceptTestQuote(now)
	expired.ExpireDate = uint32(now.Unix() - 600)
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&expired, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
}

func TestAcceptQuoteUseCase_Run_SignError(t *testing.T) {
	quoteMock := acceptTestQuote(time.Now())
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&quoteMock, nil).Once()
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("SignPegoutQuote", test.AnyCtx, acceptPegoutQuoteHash).Return("", assert.AnError).Once()

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.Error(t, err)
	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
}

func TestAcceptQuoteUseCase_S15_5_AcceptSpamLeavesLiquidityUnchanged(t *testing.T) {
	quoteHashes := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}
	quoteMock := acceptTestQuote(time.Now())
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	for _, hash := range quoteHashes {
		q := quoteMock
		quoteRepository.On("GetQuote", test.AnyCtx, hash).Return(&q, nil).Once()
		quoteRepository.On("GetRetainedQuote", test.AnyCtx, hash).Return(nil, nil).Once()
	}
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx,
		quote.PegoutStateClaimPending, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPegoutQuote{}, nil).Twice()

	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.On("GetAddress").Return("0xabcd01")

	signLp := new(mocks.ProviderMock)
	for _, hash := range quoteHashes {
		signLp.On("SignPegoutQuote", test.AnyCtx, hash).Return(acceptPegoutQuoteHashSignature, nil).Once()
	}

	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10_000), nil).Twice()
	liquidityLp := dataproviders.NewLocalLiquidityProvider(
		nil, quoteRepository, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{},
	)

	before, err := liquidityLp.AvailablePegoutLiquidity(context.Background())
	require.NoError(t, err)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, signLp,
		new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	for _, hash := range quoteHashes {
		result, runErr := useCase.Run(context.Background(), hash, "")
		require.NoError(t, runErr)
		assert.Equal(t, acceptPegoutQuoteHashSignature, result.Signature)
	}

	after, err := liquidityLp.AvailablePegoutLiquidity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, before, after)
	assert.Equal(t, entities.NewWei(10_000), after)

	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
	signLp.AssertNotCalled(t, "HasPegoutLiquidity")
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
}

func TestAcceptQuoteUseCase_Run_GetQuoteError(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, assert.AnError).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.AcceptPegoutQuoteId))
	quoteRepository.AssertExpectations(t)
	quoteRepository.AssertNotCalled(t, "GetRetainedQuote", mock.Anything, mock.Anything)
	lp.AssertNotCalled(t, "SignPegoutQuote", mock.Anything, mock.Anything)
}

func TestAcceptQuoteUseCase_Run_GetRetainedQuoteError(t *testing.T) {
	quoteMock := acceptTestQuote(time.Now())
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&quoteMock, nil).Once()
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, assert.AnError).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp, new(mocks.TrustedAccountRepositoryMock), crypto.Keccak256)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.AcceptPegoutQuoteId))
	quoteRepository.AssertExpectations(t)
	lp.AssertNotCalled(t, "SignPegoutQuote", mock.Anything, mock.Anything)
}

func acceptTrustedAccountSetup(t *testing.T) (quote.PegoutQuote, [32]byte, *entities.Signed[liquidity_provider.TrustedAccountDetails]) {
	t.Helper()
	trustedAccount := liquidity_provider.TrustedAccountDetails{
		Address:       ownerAccountAddress,
		BtcLockingCap: entities.NewWei(100000),
	}
	trustedAccountBytes, err := json.Marshal(trustedAccount)
	require.NoError(t, err)
	return acceptTestQuote(time.Now()),
		utils.To32Bytes(hexutil.MustDecode(utils.Prepend0x(acceptPegoutQuoteEip712Hash))),
		&entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value:     trustedAccount,
			Signature: "d1a9fe0de659875bc75252e6f5a73529ed6a5d88c9d97853ebf2ccc6e3080ecc423eee543470a80d373f1abb3a4f746264b47dda53252ddfc5d65989c1af34401c",
			Hash:      hex.EncodeToString(crypto.Keccak256(trustedAccountBytes)),
		}
}

func TestAcceptQuoteUseCase_Run_TrustedAccount(t *testing.T) {
	quoteMock, eip712Hash, signedTrustedAccount := acceptTrustedAccountSetup(t)
	f := newAcceptTrustedAccountFixture(&quoteMock, true)
	f.pegoutContract.EXPECT().HashPegoutQuoteEIP712(quoteMock).Return(eip712Hash, nil).Once()
	f.trustedAccountRepository.On("GetTrustedAccount", test.AnyCtx, strings.ToLower(ownerAccountAddress)).Return(signedTrustedAccount, nil).Once()
	f.quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, nil).Once()
	f.pegoutContract.On("GetAddress").Return("0xabcd01").Once()
	f.lp.On("SignPegoutQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(acceptPegoutQuoteHashSignature, nil).Once()

	result, err := f.useCase.Run(context.Background(), acceptPegoutQuoteHash, acceptPegoutQuoteHashSignature)

	require.NoError(t, err)
	assert.Equal(t, acceptPegoutQuoteHashSignature, result.Signature)
	assert.Equal(t, "0xabcd01", result.DepositAddress)
	f.assertExpectations(t)
}

func TestAcceptQuoteUseCase_Run_UntrustedAccount(t *testing.T) {
	quoteMock, eip712Hash, _ := acceptTrustedAccountSetup(t)
	f := newAcceptTrustedAccountFixture(&quoteMock, true)
	f.pegoutContract.EXPECT().HashPegoutQuoteEIP712(quoteMock).Return(eip712Hash, nil).Once()
	f.trustedAccountRepository.On("GetTrustedAccount", test.AnyCtx, strings.ToLower(ownerAccountAddress)).
		Return(nil, liquidity_provider.TrustedAccountNotFoundError).Once()

	result, err := f.useCase.Run(context.Background(), acceptPegoutQuoteHash, acceptPegoutQuoteHashSignature)

	assert.Empty(t, result)
	require.ErrorIs(t, err, liquidity_provider.TrustedAccountNotFoundError)
	f.assertExpectations(t)
	f.assertNotSigned(t)
}

func TestAcceptQuoteUseCase_Run_TamperedTrustedAccount(t *testing.T) {
	quoteMock, eip712Hash, signedTrustedAccount := acceptTrustedAccountSetup(t)
	f := newAcceptTrustedAccountFixture(&quoteMock, false)
	f.pegoutContract.EXPECT().HashPegoutQuoteEIP712(quoteMock).Return(eip712Hash, nil).Once()
	f.trustedAccountRepository.On("GetTrustedAccount", test.AnyCtx, strings.ToLower(ownerAccountAddress)).Return(signedTrustedAccount, nil).Once()

	result, err := f.useCase.Run(context.Background(), acceptPegoutQuoteHash, acceptPegoutQuoteHashSignature)

	assert.Empty(t, result)
	require.ErrorIs(t, err, liquidity_provider.TamperedTrustedAccountError)
	f.assertExpectations(t)
	f.assertNotSigned(t)
}

func TestAcceptQuoteUseCase_Run_InvalidSignature(t *testing.T) {
	quoteMock, eip712Hash, _ := acceptTrustedAccountSetup(t)
	f := newAcceptTrustedAccountFixture(&quoteMock, true)
	f.pegoutContract.EXPECT().HashPegoutQuoteEIP712(quoteMock).Return(eip712Hash, nil).Once()
	invalidSignature := "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1d"

	result, err := f.useCase.Run(context.Background(), acceptPegoutQuoteHash, invalidSignature)

	assert.Empty(t, result)
	require.ErrorContains(t, err, "recovery failed")
	f.assertExpectations(t)
	f.assertNotSigned(t)
	f.trustedAccountRepository.AssertNotCalled(t, "GetTrustedAccount", mock.Anything, mock.Anything)
}

func TestAcceptQuoteUseCase_Run_SignatureHashError(t *testing.T) {
	quoteMock, _, _ := acceptTrustedAccountSetup(t)
	f := newAcceptTrustedAccountFixture(&quoteMock, true)
	f.pegoutContract.EXPECT().HashPegoutQuoteEIP712(quoteMock).Return([32]byte{}, assert.AnError).Once()

	result, err := f.useCase.Run(context.Background(), acceptPegoutQuoteHash, acceptPegoutQuoteHashSignature)

	assert.Empty(t, result)
	require.ErrorIs(t, err, assert.AnError)
	f.assertExpectations(t)
	f.assertNotSigned(t)
}

type acceptTrustedAccountFixture struct {
	quoteRepository          *mocks.PegoutQuoteRepositoryMock
	pegoutContract           *mocks.PegoutContractMock
	lp                       *mocks.ProviderMock
	trustedAccountRepository *mocks.TrustedAccountRepositoryMock
	useCase                  *pegout.AcceptQuoteUseCase
}

func newAcceptTrustedAccountFixture(quoteMock *quote.PegoutQuote, validTrustedAccountSignature bool) acceptTrustedAccountFixture {
	f := acceptTrustedAccountFixture{
		quoteRepository:          new(mocks.PegoutQuoteRepositoryMock),
		pegoutContract:           new(mocks.PegoutContractMock),
		lp:                       new(mocks.ProviderMock),
		trustedAccountRepository: new(mocks.TrustedAccountRepositoryMock),
	}
	signer := &mocks.SignerMock{}
	signer.On("Validate", mock.Anything, mock.Anything).Return(validTrustedAccountSignature)
	f.quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(quoteMock, nil).Once()
	f.pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	f.lp.On("GetSigner").Return(signer)
	f.useCase = pegout.NewAcceptQuoteUseCase(
		f.quoteRepository,
		blockchain.RskContracts{PegOut: f.pegoutContract},
		f.lp,
		f.trustedAccountRepository,
		crypto.Keccak256,
	)
	return f
}

func (f acceptTrustedAccountFixture) assertExpectations(t *testing.T) {
	f.quoteRepository.AssertExpectations(t)
	f.pegoutContract.AssertExpectations(t)
	f.lp.AssertExpectations(t)
	f.trustedAccountRepository.AssertExpectations(t)
}

func (f acceptTrustedAccountFixture) assertNotSigned(t *testing.T) {
	f.quoteRepository.AssertNotCalled(t, "GetRetainedQuote", mock.Anything, mock.Anything)
	f.lp.AssertNotCalled(t, "SignPegoutQuote", mock.Anything, mock.Anything)
}
