package pegin_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	claimRskAddress            = "0x79568c2989232dCa1840087D73d403602364c0D4"
	claimBtcAddress            = "2N4bCrVMPu6ixHvxRS7qXcpRkdFh6WW4wEz"
	claimDepositTxID           = "9c2d9d5ab8b7e3fe1c2f2e70ad0d77f0e3c5a3e8f1b7a4e6d2c0f9e8b7a6c5d4"
	claimPegInID               = "3f32e6864ca536da70ae7ada3a0cb9c2a11b97c8a0ef9458d0f2eee6756ff8f0"
	claimTxHash                = "0x7a3c5e8f9b1d2a4c6e8f0a2b4c6d8e0f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d"
	claimClaimer               = "0x9D93929A9099be4355fC2389FbF253982F9dF47c"
	claimGasLimit              = uint64(300000)
	claimRequiredConfirmations = uint64(6)
)

var (
	claimRawTx        = []byte{0x02, 0x00, 0x00, 0x00, 0x01, 0x7f}
	claimBlockHash    = [32]byte{0x11, 0x22}
	claimMerkleBranch = blockchain.MerkleBranch{Hashes: [][32]byte{{0x33}, {0x44}}, Path: big.NewInt(5)}
)

func mustWei(t *testing.T, value string) *entities.Wei {
	t.Helper()
	parsed, ok := new(big.Int).SetString(value, 10)
	require.True(t, ok)
	return entities.NewBigWei(parsed)
}

func claimDepositAmount() *entities.Wei { return entities.NewUWei(5000000000000000000) }

func claimFee() *entities.Wei { return entities.NewWei(51000000000000000) }

func claimGasFee() *entities.Wei { return entities.NewWei(100000000000000) }

func claimValue() *entities.Wei { return entities.NewUWei(4948900000000000000) }

func claimGasPrice() *entities.Wei { return entities.NewWei(60000000) }

func claimRequiredBalance() *entities.Wei { return entities.NewUWei(4948918000000000000) }

func claimConfiguration() blockchain.PegInConfiguration {
	return blockchain.PegInConfiguration{
		FixedFee:      entities.NewWei(1000000000000000),
		PercentageFee: 100,
		MinAmount:     entities.NewWei(5000000000000000),
		MaxAmount:     entities.NewBigWei(new(big.Int).Mul(big.NewInt(100), big.NewInt(1000000000000000000))),
		RegistrantFee: entities.NewWei(100000000000000),
		ClaimGasFee:   claimGasFee(),
		ConfirmationTiers: []blockchain.ConfirmationTier{
			{MaxAmount: entities.NewUWei(10000000000000000000), Confirmations: claimRequiredConfirmations},
		},
	}
}

func claimParams() blockchain.RequestPegInParams {
	return blockchain.RequestPegInParams{
		RskAddress:         claimRskAddress,
		BtcRawTx:           claimRawTx,
		BtcBlockHeaderHash: claimBlockHash,
		MerkleBranchPath:   claimMerkleBranch.Path,
		MerkleBranchHashes: claimMerkleBranch.Hashes,
	}
}

func claimResult(callSuccess bool, eventPegInID string) blockchain.RequestPegInResult {
	return blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{TransactionHash: claimTxHash, BlockNumber: 123},
		Event: blockchain.PegInRequestedEvent{
			PegInID:     eventPegInID,
			Claimer:     claimClaimer,
			RskAddress:  claimRskAddress,
			Amount:      claimDepositAmount(),
			NetToUser:   claimValue(),
			CallSuccess: callSuccess,
		},
	}
}

func expectedClaim(state rootstock.PegInClaimState, stored rootstock.PegInClaim) rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:  claimRskAddress,
		DepositTxID: claimDepositTxID,
		BtcAddress:  claimBtcAddress,
		PegInID:     claimPegInID,
		State:       state,
		Amount:      claimDepositAmount(),
		Fee:         claimFee(),
		GasFee:      claimGasFee(),
		Value:       claimValue(),
		CreatedAt:   stored.CreatedAt,
		UpdatedAt:   stored.UpdatedAt,
	}
}

type claimPegInFixtures struct {
	claims         *mocks.PegInClaimRepositoryMock
	pegIn          *mocks.PeginContractMock
	configurations *mocks.FlyoverConfigurationsContractMock
	pauseRegistry  *mocks.PauseRegistryContractMock
	btcRpc         *mocks.BtcRpcMock
	rskRpc         *mocks.RootstockRpcServerMock
	wallet         *mocks.RskSignerWalletMock
	liquidityMutex *sync.Mutex
	walletMutex    *sync.Mutex
	useCase        *pegin.ClaimPegInUseCase
}

func newClaimPegInFixtures() *claimPegInFixtures {
	f := &claimPegInFixtures{
		claims:         &mocks.PegInClaimRepositoryMock{},
		pegIn:          &mocks.PeginContractMock{},
		configurations: &mocks.FlyoverConfigurationsContractMock{},
		pauseRegistry:  &mocks.PauseRegistryContractMock{},
		btcRpc:         &mocks.BtcRpcMock{},
		rskRpc:         &mocks.RootstockRpcServerMock{},
		wallet:         &mocks.RskSignerWalletMock{},
		liquidityMutex: &sync.Mutex{},
		walletMutex:    &sync.Mutex{},
	}
	f.useCase = pegin.NewClaimPegInUseCase(
		f.claims,
		blockchain.RskContracts{PegIn: f.pegIn, FlyoverConfigurations: f.configurations, PauseRegistry: f.pauseRegistry},
		blockchain.Rpc{Btc: f.btcRpc, Rsk: f.rskRpc},
		f.wallet,
		crypto.Keccak256,
		f.liquidityMutex,
		f.walletMutex,
	)
	return f
}

func (f *claimPegInFixtures) run() error {
	return f.runFor(claimRskAddress, claimDepositTxID)
}

func (f *claimPegInFixtures) runFor(rskAddress, depositTxID string) error {
	watch := rootstock.PegInWatch{RskAddress: rskAddress, BtcAddress: claimBtcAddress, State: rootstock.PegInWatchImported}
	return f.useCase.Run(context.Background(), watch, depositTxID)
}

func (f *claimPegInFixtures) assertExpectations(t *testing.T) {
	f.claims.AssertExpectations(t)
	f.pegIn.AssertExpectations(t)
	f.configurations.AssertExpectations(t)
	f.pauseRegistry.AssertExpectations(t)
	f.btcRpc.AssertExpectations(t)
	f.rskRpc.AssertExpectations(t)
	f.wallet.AssertExpectations(t)
}

func (f *claimPegInFixtures) expectNoStoredClaim() {
	f.claims.EXPECT().Get(mock.Anything, claimRskAddress, claimDepositTxID).Return(nil, nil).Once()
}

func (f *claimPegInFixtures) expectDepositInfo(amount *entities.Wei, confirmations uint64) {
	f.btcRpc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: confirmations,
		Outputs: map[string][]*entities.Wei{
			claimBtcAddress:                       {amount, entities.NewWei(7)},
			"2NEbgUn5XcHwUnVVQUZ2hzVYvLzfqPcBdAG": {entities.NewWei(9)},
		},
	}, nil).Once()
}

func (f *claimPegInFixtures) expectConfirmationTier(amount *entities.Wei) {
	f.configurations.EXPECT().GetRequiredPegInBtcConfirmations(amount).Return(claimRequiredConfirmations, nil).Once()
}

func (f *claimPegInFixtures) expectConfirmedDeposit(amount *entities.Wei) {
	f.expectNoStoredClaim()
	f.expectDepositInfo(amount, claimRequiredConfirmations)
	f.expectConfirmationTier(amount)
	f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
	f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(claimRawTx, nil).Once()
}

func (f *claimPegInFixtures) expectSpvProof() {
	f.btcRpc.On("GetTransactionBlockInfo", claimDepositTxID).
		Return(blockchain.BitcoinBlockInformation{Hash: claimBlockHash, Height: big.NewInt(800)}, nil).Once()
	f.btcRpc.On("BuildMerkleBranch", claimDepositTxID).Return(claimMerkleBranch, nil).Once()
}

func (f *claimPegInFixtures) expectPreparedClaim() {
	f.expectConfirmedDeposit(claimDepositAmount())
	f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
	f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(claimFee(), nil).Once()
	f.expectSpvProof()
}

func (f *claimPegInFixtures) expectDryRun(value *entities.Wei, err error) {
	f.pegIn.EXPECT().SimulateRequestPegIn(blockchain.NewTransactionConfig(value, 0, nil), claimParams()).Return(err).Once()
}

func (f *claimPegInFixtures) expectFunding(value, balance *entities.Wei) {
	f.pegIn.EXPECT().EstimateRequestPegInGas(blockchain.NewTransactionConfig(value, 0, nil), claimParams()).
		Return(claimGasLimit, nil).Once()
	f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(claimGasPrice(), nil).Once()
	f.wallet.EXPECT().GetBalance(mock.Anything).Return(balance, nil).Once()
}

func (f *claimPegInFixtures) expectAdd(err error) *rootstock.PegInClaim {
	added := &rootstock.PegInClaim{}
	f.claims.EXPECT().Add(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, claim rootstock.PegInClaim) error {
			*added = claim
			return err
		}).Once()
	return added
}

func (f *claimPegInFixtures) expectUpdate(err error) *rootstock.PegInClaim {
	updated := &rootstock.PegInClaim{}
	f.claims.EXPECT().Update(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, claim rootstock.PegInClaim) error {
			*updated = claim
			return err
		}).Once()
	return updated
}

func (f *claimPegInFixtures) expectSend(result blockchain.RequestPegInResult, err error) {
	f.pegIn.EXPECT().RequestPegIn(blockchain.NewTransactionConfig(claimValue(), claimGasLimit, claimGasPrice()), claimParams()).
		Return(result, err).Once()
}

func (f *claimPegInFixtures) expectSubmittedClaim() (added *rootstock.PegInClaim) {
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), nil)
	f.expectFunding(claimValue(), claimRequiredBalance())
	return f.expectAdd(nil)
}

func TestClaimPegInUseCase_Run_ClaimsAConfirmedDeposit(t *testing.T) {
	f := newClaimPegInFixtures()
	added := f.expectSubmittedClaim()
	f.expectSend(claimResult(true, claimPegInID), nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInClaimed(claimRskAddress, claimDepositTxID, claimTxHash, claimPegInID))()

	require.NoError(t, f.run())

	assert.Equal(t, expectedClaim(rootstock.PegInClaimSubmitting, *added), *added)
	assert.WithinDuration(t, time.Now(), added.CreatedAt, time.Minute)
	claimed := expectedClaim(rootstock.PegInClaimClaimed, *updated)
	claimed.TxHash = claimTxHash
	claimed.CallSuccess = true
	assert.Equal(t, claimed, *updated)
	assert.Equal(t, added.CreatedAt, updated.CreatedAt)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SendsOnlyUpToTheMaximum(t *testing.T) {
	f := newClaimPegInFixtures()
	maxAmount := mustWei(t, "100000000000000000000")
	amount := mustWei(t, "100000000010000000000")
	fee := entities.NewUWei(1001000000000000000)
	value := mustWei(t, "98999000000000000000")
	configuration := claimConfiguration()
	configuration.MaxAmount = maxAmount
	configuration.ClaimGasFee = entities.NewWei(0)
	f.expectConfirmedDeposit(amount)
	f.configurations.EXPECT().GetPegInConfiguration().Return(configuration, nil).Once()
	f.configurations.EXPECT().CalculatePegInFee(maxAmount).Return(fee, nil).Once()
	f.expectSpvProof()
	f.expectDryRun(value, blockchain.ContractPausedError)

	require.NoError(t, f.run())
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SkipsADepositWithAStoredClaim(t *testing.T) {
	for _, state := range []rootstock.PegInClaimState{
		rootstock.PegInClaimSubmitting,
		rootstock.PegInClaimClaimed,
		rootstock.PegInClaimRaceLost,
	} {
		t.Run(string(state), func(t *testing.T) {
			log.SetLevel(log.DebugLevel)
			f := newClaimPegInFixtures()
			stored := rootstock.PegInClaim{RskAddress: claimRskAddress, DepositTxID: claimDepositTxID, State: state}
			f.claims.EXPECT().Get(mock.Anything, claimRskAddress, claimDepositTxID).Return(&stored, nil).Once()
			defer test.AssertLogContains(t, pegin.LogPegInClaimAlreadyStored(claimRskAddress, claimDepositTxID, state))()

			require.NoError(t, f.run())
			f.btcRpc.AssertNotCalled(t, "GetTransactionInfo", claimDepositTxID)
			f.claims.AssertExpectations(t)
		})
	}
}

func TestClaimPegInUseCase_Run_RetriesARetryableFailureKeepingItsCreationTime(t *testing.T) {
	f := newClaimPegInFixtures()
	firstAttempt := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	stored := rootstock.PegInClaim{
		RskAddress:  claimRskAddress,
		DepositTxID: claimDepositTxID,
		State:       rootstock.PegInClaimRetryableFailure,
		CreatedAt:   firstAttempt,
	}
	f.claims.EXPECT().Get(mock.Anything, claimRskAddress, claimDepositTxID).Return(&stored, nil).Once()
	f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations)
	f.expectConfirmationTier(claimDepositAmount())
	f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
	f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(claimRawTx, nil).Once()
	f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
	f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(claimFee(), nil).Once()
	f.expectSpvProof()
	f.expectDryRun(claimValue(), nil)
	f.expectFunding(claimValue(), claimRequiredBalance())
	added := f.expectAdd(nil)
	f.expectSend(claimResult(true, claimPegInID), nil)
	f.expectUpdate(nil)

	require.NoError(t, f.run())
	assert.Equal(t, firstAttempt, added.CreatedAt)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_WaitsForTheConfirmationTier(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newClaimPegInFixtures()
	f.expectNoStoredClaim()
	f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations-1)
	f.expectConfirmationTier(claimDepositAmount())
	defer test.AssertLogContains(t, pegin.LogPegInClaimWaitingForConfirmations(
		claimRskAddress, claimDepositTxID, claimRequiredConfirmations-1, claimRequiredConfirmations,
	))()

	require.NoError(t, f.run())
	f.pauseRegistry.AssertNotCalled(t, "PauseLevel")
	f.btcRpc.AssertNotCalled(t, "GetRawTransaction", claimDepositTxID)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_FailsForADepositWithoutAnOutputToTheWatchedAddress(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectNoStoredClaim()
	f.btcRpc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: claimRequiredConfirmations,
		Outputs:       map[string][]*entities.Wei{"2NEbgUn5XcHwUnVVQUZ2hzVYvLzfqPcBdAG": {entities.NewWei(9)}},
	}, nil).Once()

	err := f.run()
	require.ErrorIs(t, err, pegin.DepositOutputNotFoundError)
	require.ErrorContains(t, err, string(usecases.ClaimPegInId))
	f.configurations.AssertNotCalled(t, "GetRequiredPegInBtcConfirmations", mock.Anything)
}

func TestClaimPegInUseCase_Run_SkipsUnderAHardPause(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newClaimPegInFixtures()
	f.expectNoStoredClaim()
	f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations)
	f.expectConfirmationTier(claimDepositAmount())
	f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelHard, nil).Once()
	defer test.AssertLogContains(t, pegin.LogPegInClaimHardPaused(claimRskAddress, claimDepositTxID))()

	require.NoError(t, f.run())
	f.btcRpc.AssertNotCalled(t, "GetRawTransaction", claimDepositTxID)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_LetsTheDryRunDecideUnderASoftPause(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectNoStoredClaim()
	f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations)
	f.expectConfirmationTier(claimDepositAmount())
	f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelSoft, nil).Once()
	f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(claimRawTx, nil).Once()
	f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
	f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(claimFee(), nil).Once()
	f.expectSpvProof()
	f.expectDryRun(claimValue(), blockchain.ContractPausedError)

	require.NoError(t, f.run())
	f.claims.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_RejectsAnOversizedDeposit(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectNoStoredClaim()
	f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations)
	f.expectConfirmationTier(claimDepositAmount())
	f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
	f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(make([]byte, pegin.MaxPeginDepositTxSize+1), nil).Once()

	err := f.run()
	require.ErrorIs(t, err, pegin.DataCapExceededError)
	require.ErrorContains(t, err, string(usecases.ClaimPegInId))
	f.configurations.AssertNotCalled(t, "GetPegInConfiguration")
	f.pegIn.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything, mock.Anything)
	f.claims.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_Run_AcceptsADepositAtTheSizeLimit(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectNoStoredClaim()
	f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations)
	f.expectConfirmationTier(claimDepositAmount())
	f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
	f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(make([]byte, pegin.MaxPeginDepositTxSize), nil).Once()
	configuration := claimConfiguration()
	configuration.MinAmount = new(entities.Wei).Add(claimDepositAmount(), entities.NewWei(1))
	f.configurations.EXPECT().GetPegInConfiguration().Return(configuration, nil).Once()

	require.NoError(t, f.run())
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SkipsADepositBelowTheMinimum(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newClaimPegInFixtures()
	f.expectConfirmedDeposit(claimDepositAmount())
	configuration := claimConfiguration()
	configuration.MinAmount = new(entities.Wei).Add(claimDepositAmount(), entities.NewWei(1))
	f.configurations.EXPECT().GetPegInConfiguration().Return(configuration, nil).Once()
	defer test.AssertLogContains(t, pegin.LogPegInClaimBelowMinimum(
		claimRskAddress, claimDepositTxID, claimDepositAmount(), configuration.MinAmount,
	))()

	require.NoError(t, f.run())
	f.configurations.AssertNotCalled(t, "CalculatePegInFee", mock.Anything)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SkipsWhenTheFeesExceedTheAmount(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newClaimPegInFixtures()
	f.expectConfirmedDeposit(claimDepositAmount())
	f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
	fee := new(entities.Wei).Sub(claimDepositAmount(), claimGasFee())
	fee.Add(fee, entities.NewWei(1))
	f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(fee, nil).Once()
	defer test.AssertLogContains(t, pegin.LogPegInClaimDeductionAboveAmount(
		claimRskAddress, claimDepositTxID, claimDepositAmount(), new(entities.Wei).Add(claimDepositAmount(), entities.NewWei(1)),
	))()

	require.NoError(t, f.run())
	f.btcRpc.AssertNotCalled(t, "GetTransactionBlockInfo", claimDepositTxID)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_ClaimsADepositTheFeesConsumeExactly(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectConfirmedDeposit(claimDepositAmount())
	f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
	fee := new(entities.Wei).Sub(claimDepositAmount(), claimGasFee())
	f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(fee, nil).Once()
	f.expectSpvProof()
	f.pegIn.EXPECT().SimulateRequestPegIn(mock.MatchedBy(func(txConfig blockchain.TransactionConfig) bool {
		return txConfig.Value.Cmp(entities.NewWei(0)) == 0 && txConfig.GasLimit == nil
	}), claimParams()).Return(blockchain.ContractPausedError).Once()

	require.NoError(t, f.run())
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_RecordsARaceLostInTheDryRun(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), blockchain.PegInAlreadyProcessedError)
	added := f.expectAdd(nil)
	defer test.AssertLogContains(t, pegin.LogPegInClaimRaceLost(claimRskAddress, claimDepositTxID))()

	require.NoError(t, f.run())
	assert.Equal(t, expectedClaim(rootstock.PegInClaimRaceLost, *added), *added)
	f.pegIn.AssertNotCalled(t, "EstimateRequestPegInGas", mock.Anything, mock.Anything)
	f.pegIn.AssertNotCalled(t, "RequestPegIn", mock.Anything, mock.Anything)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_IgnoresARaceLostClaimStoredMeanwhile(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), blockchain.PegInAlreadyProcessedError)
	f.expectAdd(rootstock.PegInClaimExistsError)

	require.NoError(t, f.run())
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_ReturnsTheRaceLostStoreError(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), blockchain.PegInAlreadyProcessedError)
	f.expectAdd(assert.AnError)

	err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ClaimPegInId))
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SkipsWhileTheDryRunSaysNotYet(t *testing.T) {
	for _, dryRunErr := range []error{
		blockchain.PegInBelowMinimumError,
		blockchain.DepositBelowBridgeMinimumError,
		blockchain.InsufficientConfirmationsError,
		blockchain.UnableToGetConfirmationsError,
		blockchain.ContractPausedError,
	} {
		t.Run(dryRunErr.Error(), func(t *testing.T) {
			log.SetLevel(log.DebugLevel)
			f := newClaimPegInFixtures()
			f.expectPreparedClaim()
			f.expectDryRun(claimValue(), dryRunErr)
			defer test.AssertLogContains(t, pegin.LogPegInClaimNotClaimable(claimRskAddress, claimDepositTxID, dryRunErr))()

			require.NoError(t, f.run())
			f.claims.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
			f.assertExpectations(t)
		})
	}
}

func TestClaimPegInUseCase_Run_ReturnsAnUnexpectedDryRunError(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), blockchain.IncorrectFrontingError)

	err := f.run()
	require.ErrorIs(t, err, blockchain.IncorrectFrontingError)
	require.ErrorContains(t, err, string(usecases.ClaimPegInId))
	f.claims.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SkipsWhenTheWalletCantCoverTheClaim(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), nil)
	balance := new(entities.Wei).Sub(claimRequiredBalance(), entities.NewWei(1))
	f.expectFunding(claimValue(), balance)
	defer test.AssertLogContains(t, pegin.LogPegInClaimInsufficientBalance(
		claimRskAddress, claimDepositTxID, balance, claimRequiredBalance(),
	))()

	require.NoError(t, f.run())
	f.claims.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	f.pegIn.AssertNotCalled(t, "RequestPegIn", mock.Anything, mock.Anything)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_SkipsAClaimStoredConcurrently(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), nil)
	f.expectFunding(claimValue(), claimRequiredBalance())
	f.expectAdd(rootstock.PegInClaimExistsError)

	require.NoError(t, f.run())
	f.pegIn.AssertNotCalled(t, "RequestPegIn", mock.Anything, mock.Anything)
	f.assertExpectations(t)
}

func assertMutexHeld(t *testing.T, mutex *sync.Mutex, message string) {
	acquired := mutex.TryLock()
	if acquired {
		mutex.Unlock()
	}
	assert.False(t, acquired, message)
}

func TestClaimPegInUseCase_Run_HoldsTheLiquidityMutexFromTheBalanceCheckToTheSend(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectPreparedClaim()
	f.expectDryRun(claimValue(), nil)
	f.pegIn.EXPECT().EstimateRequestPegInGas(mock.Anything, mock.Anything).Return(claimGasLimit, nil).Once()
	f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(claimGasPrice(), nil).Once()
	f.wallet.EXPECT().GetBalance(mock.Anything).RunAndReturn(func(context.Context) (*entities.Wei, error) {
		assertMutexHeld(t, f.liquidityMutex, "the balance check runs under the liquidity mutex")
		return claimRequiredBalance(), nil
	}).Once()
	f.expectAdd(nil)
	f.pegIn.EXPECT().RequestPegIn(mock.Anything, mock.Anything).
		RunAndReturn(func(blockchain.TransactionConfig, blockchain.RequestPegInParams) (blockchain.RequestPegInResult, error) {
			assertMutexHeld(t, f.liquidityMutex, "the send runs under the liquidity mutex")
			assertMutexHeld(t, f.walletMutex, "the send runs under the wallet mutex")
			return claimResult(true, claimPegInID), nil
		}).Once()
	f.expectUpdate(nil)

	require.NoError(t, f.run())
	assert.True(t, f.liquidityMutex.TryLock(), "the liquidity mutex is released after the claim")
	assert.True(t, f.walletMutex.TryLock(), "the wallet mutex is released after the claim")
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_RecordsARaceLostAtMining(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	reverted := blockchain.RequestPegInResult{Receipt: blockchain.TransactionReceipt{TransactionHash: claimTxHash}}
	f.expectSend(reverted, errors.Join(blockchain.TxFailedError, blockchain.PegInAlreadyProcessedError))
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInClaimRaceLost(claimRskAddress, claimDepositTxID))()

	require.NoError(t, f.run())
	raceLost := expectedClaim(rootstock.PegInClaimRaceLost, *updated)
	raceLost.TxHash = claimTxHash
	assert.Equal(t, raceLost, *updated)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_MarksAClaimThatWasNotBroadcastAsRetryable(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	f.expectSend(blockchain.RequestPegInResult{}, assert.AnError)
	updated := f.expectUpdate(nil)

	err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ClaimPegInId))
	assert.Equal(t, expectedClaim(rootstock.PegInClaimRetryableFailure, *updated), *updated)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_MarksAMinedRevertAsRetryable(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	reverted := blockchain.RequestPegInResult{Receipt: blockchain.TransactionReceipt{TransactionHash: claimTxHash}}
	f.expectSend(reverted, errors.Join(blockchain.TxFailedError, blockchain.IncorrectFrontingError))
	updated := f.expectUpdate(nil)

	err := f.run()
	require.ErrorIs(t, err, blockchain.TxFailedError)
	assert.Equal(t, expectedClaim(rootstock.PegInClaimRetryableFailure, *updated), *updated)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_KeepsAClaimWithAnUnknownOutcomeSubmitting(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	broadcast := blockchain.RequestPegInResult{Receipt: blockchain.TransactionReceipt{TransactionHash: claimTxHash}}
	f.expectSend(broadcast, assert.AnError)
	updated := f.expectUpdate(nil)

	err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	submitting := expectedClaim(rootstock.PegInClaimSubmitting, *updated)
	submitting.TxHash = claimTxHash
	assert.Equal(t, submitting, *updated)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_ReturnsTheUpdateErrorAfterTheSend(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	f.expectSend(claimResult(true, claimPegInID), nil)
	f.expectUpdate(assert.AnError)

	err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ClaimPegInId))
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_KeepsTheLocalPegInIdWhenTheEventDiffers(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	const eventPegInID = "00000000000000000000000000000000000000000000000000000000000000aa"
	f.expectSend(claimResult(true, eventPegInID), nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInClaimIdMismatch(claimDepositTxID, claimPegInID, eventPegInID))()

	require.NoError(t, f.run())
	assert.Equal(t, claimPegInID, updated.PegInID)
	assert.Equal(t, rootstock.PegInClaimClaimed, updated.State)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_RecordsAFailedUserPayment(t *testing.T) {
	f := newClaimPegInFixtures()
	f.expectSubmittedClaim()
	f.expectSend(claimResult(false, claimPegInID), nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInClaimUserPaymentFailed(claimRskAddress, claimDepositTxID, claimTxHash))()

	require.NoError(t, f.run())
	claimed := expectedClaim(rootstock.PegInClaimClaimed, *updated)
	claimed.TxHash = claimTxHash
	assert.Equal(t, claimed, *updated)
	f.assertExpectations(t)
}

func TestClaimPegInUseCase_Run_FailsWhenThePegInIdInputsAreNotHex(t *testing.T) {
	cases := []struct {
		name        string
		rskAddress  string
		depositTxID string
	}{
		{name: "rsk address", rskAddress: "0xnothex", depositTxID: claimDepositTxID},
		{name: "deposit txid", rskAddress: claimRskAddress, depositTxID: "nothex"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f := newClaimPegInFixtures()
			f.claims.EXPECT().Get(mock.Anything, testCase.rskAddress, testCase.depositTxID).Return(nil, nil).Once()
			f.btcRpc.On("GetTransactionInfo", testCase.depositTxID).Return(blockchain.BitcoinTransactionInformation{
				Hash:          testCase.depositTxID,
				Confirmations: claimRequiredConfirmations,
				Outputs:       map[string][]*entities.Wei{claimBtcAddress: {claimDepositAmount()}},
			}, nil).Once()
			f.expectConfirmationTier(claimDepositAmount())
			f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
			f.btcRpc.On("GetRawTransaction", testCase.depositTxID).Return(claimRawTx, nil).Once()
			f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
			f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(claimFee(), nil).Once()
			f.btcRpc.On("GetTransactionBlockInfo", testCase.depositTxID).
				Return(blockchain.BitcoinBlockInformation{Hash: claimBlockHash}, nil).Once()
			f.btcRpc.On("BuildMerkleBranch", testCase.depositTxID).Return(claimMerkleBranch, nil).Once()

			err := f.runFor(testCase.rskAddress, testCase.depositTxID)
			require.ErrorContains(t, err, string(usecases.ClaimPegInId))
			f.pegIn.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything, mock.Anything)
		})
	}
}

type claimStep struct {
	name    string
	succeed func(f *claimPegInFixtures)
	fail    func(f *claimPegInFixtures)
}

// nolint:funlen
func claimSteps() []claimStep {
	return []claimStep{
		{
			name:    "stored claim read",
			succeed: (*claimPegInFixtures).expectNoStoredClaim,
			fail: func(f *claimPegInFixtures) {
				f.claims.EXPECT().Get(mock.Anything, claimRskAddress, claimDepositTxID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "deposit read",
			succeed: func(f *claimPegInFixtures) { f.expectDepositInfo(claimDepositAmount(), claimRequiredConfirmations) },
			fail: func(f *claimPegInFixtures) {
				f.btcRpc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
			},
		},
		{
			name:    "confirmation tier read",
			succeed: func(f *claimPegInFixtures) { f.expectConfirmationTier(claimDepositAmount()) },
			fail: func(f *claimPegInFixtures) {
				f.configurations.EXPECT().GetRequiredPegInBtcConfirmations(claimDepositAmount()).Return(0, assert.AnError).Once()
			},
		},
		{
			name: "pause level read",
			succeed: func(f *claimPegInFixtures) {
				f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.pauseRegistry.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, assert.AnError).Once()
			},
		},
		{
			name: "raw transaction read",
			succeed: func(f *claimPegInFixtures) {
				f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(claimRawTx, nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "configuration read",
			succeed: func(f *claimPegInFixtures) {
				f.configurations.EXPECT().GetPegInConfiguration().Return(claimConfiguration(), nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.configurations.EXPECT().GetPegInConfiguration().Return(blockchain.PegInConfiguration{}, assert.AnError).Once()
			},
		},
		{
			name: "fee read",
			succeed: func(f *claimPegInFixtures) {
				f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(claimFee(), nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.configurations.EXPECT().CalculatePegInFee(claimDepositAmount()).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "block read",
			succeed: func(f *claimPegInFixtures) {
				f.btcRpc.On("GetTransactionBlockInfo", claimDepositTxID).
					Return(blockchain.BitcoinBlockInformation{Hash: claimBlockHash, Height: big.NewInt(800)}, nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.btcRpc.On("GetTransactionBlockInfo", claimDepositTxID).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
			},
		},
		{
			name: "merkle branch read",
			succeed: func(f *claimPegInFixtures) {
				f.btcRpc.On("BuildMerkleBranch", claimDepositTxID).Return(claimMerkleBranch, nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.btcRpc.On("BuildMerkleBranch", claimDepositTxID).Return(blockchain.MerkleBranch{}, assert.AnError).Once()
			},
		},
		{
			name:    "dry run",
			succeed: func(f *claimPegInFixtures) { f.expectDryRun(claimValue(), nil) },
			fail:    func(f *claimPegInFixtures) { f.expectDryRun(claimValue(), assert.AnError) },
		},
		{
			name: "gas estimate",
			succeed: func(f *claimPegInFixtures) {
				f.pegIn.EXPECT().EstimateRequestPegInGas(blockchain.NewTransactionConfig(claimValue(), 0, nil), claimParams()).
					Return(claimGasLimit, nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.pegIn.EXPECT().EstimateRequestPegInGas(blockchain.NewTransactionConfig(claimValue(), 0, nil), claimParams()).
					Return(0, assert.AnError).Once()
			},
		},
		{
			name: "gas price read",
			succeed: func(f *claimPegInFixtures) {
				f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(claimGasPrice(), nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "wallet balance read",
			succeed: func(f *claimPegInFixtures) {
				f.wallet.EXPECT().GetBalance(mock.Anything).Return(claimRequiredBalance(), nil).Once()
			},
			fail: func(f *claimPegInFixtures) {
				f.wallet.EXPECT().GetBalance(mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "submitting claim insert",
			succeed: func(f *claimPegInFixtures) { f.expectAdd(nil) },
			fail:    func(f *claimPegInFixtures) { f.expectAdd(assert.AnError) },
		},
	}
}

func TestClaimPegInUseCase_Run_ReturnsTheErrorOfEachStep(t *testing.T) {
	steps := claimSteps()
	for i, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			f := newClaimPegInFixtures()
			for _, previous := range steps[:i] {
				previous.succeed(f)
			}
			step.fail(f)

			err := f.run()
			require.ErrorIs(t, err, assert.AnError)
			require.ErrorContains(t, err, string(usecases.ClaimPegInId))
			f.pegIn.AssertNotCalled(t, "RequestPegIn", mock.Anything, mock.Anything)
			f.assertExpectations(t)
		})
	}
}
