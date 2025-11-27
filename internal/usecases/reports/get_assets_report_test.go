package reports_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fixed fees per quote (not part of RequiredLiquidity but added to each quote's Total())
const (
	callFeePerQuote          = 1000000 // 0.000001 RBTC (1e6 wei)
	productFeeAmountPerQuote = 500000  // 0.0000005 RBTC (5e5 wei)
)

// additionalFeesPerQuote returns the sum of fixed fees added to each quote
func additionalFeesPerQuote() *entities.Wei {
	total := entities.NewWei(callFeePerQuote)
	total.Add(total, entities.NewWei(productFeeAmountPerQuote))
	return total
}

type ExpectedBtcCalculations struct {
	WalletBalance         *entities.Wei
	Rebalancing           *entities.Wei // Federation
	WaitingForRebalancing *entities.Wei // RSK Wallet
	InLbc                 *entities.Wei
	ReservedForUsers      *entities.Wei
	WaitingForRefund      *entities.Wei
	Total                 *entities.Wei
	Available             *entities.Wei
}

type splitQuoteValues struct {
	Value            *entities.Wei
	CallFee          *entities.Wei
	ProductFeeAmount *entities.Wei
	GasFee           *entities.Wei
}

// splitRequiredLiquidity splits the RequiredLiquidity into quote field values so the call to Total() is more realistic:
// RequiredLiquidity is composed of Value and GasFee. We split it as 5% for GasFee and 95% for Value.
// CallFee and ProductFeeAmount are set to reasonable fixed values (not part of the required liquidity calculation).
func splitRequiredLiquidity(t *testing.T, requiredLiquidity *entities.Wei) splitQuoteValues {
	callFee := entities.NewWei(callFeePerQuote)
	productFeeAmount := entities.NewWei(productFeeAmountPerQuote)

	// Calculate 5% for GasFee
	gasFeeInt := new(big.Int).Mul(requiredLiquidity.AsBigInt(), big.NewInt(5))
	gasFeeInt.Div(gasFeeInt, big.NewInt(100))
	gasFee := entities.NewBigWei(gasFeeInt)

	// Calculate 95% for Value (remainder ensures perfect precision)
	valueBigInt := new(big.Int).Sub(requiredLiquidity.AsBigInt(), gasFeeInt)
	value := entities.NewBigWei(valueBigInt)

	valueAndGas := new(entities.Wei).Add(value, gasFee)
	require.Equal(t, 0, valueAndGas.Cmp(requiredLiquidity),
		"Value + GasFee must equal RequiredLiquidity")

	return splitQuoteValues{
		Value:            value,
		CallFee:          callFee,
		ProductFeeAmount: productFeeAmount,
		GasFee:           gasFee,
	}
}

// The states are stored in the retained quotes, so the filtering is done based on those structures but the amount calculations are obtained
// by calling quote.Total(). Here we convert the retained quotes to the regular quotes to be able to call that method.
func retainedPegoutQuotesToPegoutQuotes(t *testing.T, retainedQuotes []quote.RetainedPegoutQuote) []quote.PegoutQuote {
	result := make([]quote.PegoutQuote, len(retainedQuotes))
	for i, rq := range retainedQuotes {
		split := splitRequiredLiquidity(t, rq.RequiredLiquidity)
		result[i] = quote.PegoutQuote{
			Value:            split.Value,
			CallFee:          split.CallFee,
			ProductFeeAmount: split.ProductFeeAmount,
			GasFee:           split.GasFee,
		}
	}
	return result
}

func retainedPeginQuotesToPeginQuotes(t *testing.T, retainedQuotes []quote.RetainedPeginQuote) []quote.PeginQuote {
	result := make([]quote.PeginQuote, len(retainedQuotes))
	for i, rq := range retainedQuotes {
		split := splitRequiredLiquidity(t, rq.RequiredLiquidity)
		result[i] = quote.PeginQuote{
			Value:            split.Value,
			CallFee:          split.CallFee,
			ProductFeeAmount: split.ProductFeeAmount,
			GasFee:           split.GasFee,
		}
	}
	return result
}

type ExpectedRbtcCalculations struct {
	RskWalletBalance *entities.Wei // Raw RSK wallet balance
	InRskWallet      *entities.Wei // Adjusted RSK wallet (subtracting BTC waiting for rebalancing)
	LockedInLbc      *entities.Wei
	WaitingForRefund *entities.Wei // Federation
	ReservedForUsers *entities.Wei
	Total            *entities.Wei
	Available        *entities.Wei
}

// nolint:exhaustive
func calculateExpectedBtcValues(quotes []quote.RetainedPegoutQuote, btcWalletBalance *entities.Wei) ExpectedBtcCalculations {
	expectedBtcRebalancing := entities.NewWei(0)           // BridgeTxSucceeded
	expectedBtcWaitingForRebalancing := entities.NewWei(0) // RefundPegOutSucceeded
	expectedBtcInLbc := entities.NewWei(0)                 // SendPegoutSucceeded
	expectedBtcReservedForUsers := entities.NewWei(0)      // WaitingForDeposit + WaitingForDepositConfirmations
	expectedBtcWaitingForRefund := entities.NewWei(0)      // RefundPegOutSucceeded + SendPegoutSucceeded + BridgeTxSucceeded

	// Calculate sums based on quote states
	// Each quote's Total() = RequiredLiquidity + callFeePerQuote + productFeeAmountPerQuote
	additionalFees := additionalFeesPerQuote()

	for _, q := range quotes {
		quoteTotal := new(entities.Wei).Add(q.RequiredLiquidity, additionalFees)

		switch q.State {
		case quote.PegoutStateBridgeTxSucceeded:
			expectedBtcRebalancing.Add(expectedBtcRebalancing, quoteTotal)
			expectedBtcWaitingForRefund.Add(expectedBtcWaitingForRefund, quoteTotal)
		case quote.PegoutStateRefundPegOutSucceeded:
			expectedBtcWaitingForRebalancing.Add(expectedBtcWaitingForRebalancing, quoteTotal)
			expectedBtcWaitingForRefund.Add(expectedBtcWaitingForRefund, quoteTotal)
		case quote.PegoutStateSendPegoutSucceeded:
			expectedBtcInLbc.Add(expectedBtcInLbc, quoteTotal)
			expectedBtcWaitingForRefund.Add(expectedBtcWaitingForRefund, quoteTotal)
		case quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations:
			expectedBtcReservedForUsers.Add(expectedBtcReservedForUsers, quoteTotal)
		}
	}

	expectedBtcTotal := entities.NewWei(0)
	expectedBtcTotal.Add(expectedBtcTotal, btcWalletBalance)
	expectedBtcTotal.Add(expectedBtcTotal, expectedBtcRebalancing)
	expectedBtcTotal.Add(expectedBtcTotal, expectedBtcWaitingForRebalancing)
	expectedBtcTotal.Add(expectedBtcTotal, expectedBtcInLbc)

	expectedBtcAvailable := entities.NewWei(0).Sub(btcWalletBalance, expectedBtcReservedForUsers)

	return ExpectedBtcCalculations{
		WalletBalance:         btcWalletBalance,
		Rebalancing:           expectedBtcRebalancing,
		WaitingForRebalancing: expectedBtcWaitingForRebalancing,
		InLbc:                 expectedBtcInLbc,
		ReservedForUsers:      expectedBtcReservedForUsers,
		WaitingForRefund:      expectedBtcWaitingForRefund,
		Total:                 expectedBtcTotal,
		Available:             expectedBtcAvailable,
	}
}

// nolint:exhaustive
func calculateExpectedRbtcValues(
	peginQuotes []quote.RetainedPeginQuote,
	rbtcWalletBalance *entities.Wei,
	rbtcLockedInLbc *entities.Wei,
	btcWaitingForRebalancing *entities.Wei,
) ExpectedRbtcCalculations {
	expectedRbtcWaitingForRefund := entities.NewWei(0) // CallForUserSucceeded
	expectedRbtcReservedForUsers := entities.NewWei(0) // WaitingForDeposit + WaitingForDepositConfirmations

	// Calculate sums based on pegin quote states
	// Each quote's Total() = RequiredLiquidity + callFeePerQuote + productFeeAmountPerQuote
	additionalFees := additionalFeesPerQuote()

	for _, q := range peginQuotes {
		quoteTotal := new(entities.Wei).Add(q.RequiredLiquidity, additionalFees)

		switch q.State {
		case quote.PeginStateCallForUserSucceeded:
			expectedRbtcWaitingForRefund.Add(expectedRbtcWaitingForRefund, quoteTotal)
		case quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations:
			expectedRbtcReservedForUsers.Add(expectedRbtcReservedForUsers, quoteTotal)
		}
	}

	// A part of the RBTC in the RSK wallet is a representation of BTC waiting to be sent to the bridge for rebalancing
	expectedRbtcInRskWallet := entities.NewWei(0).Sub(rbtcWalletBalance, btcWaitingForRebalancing)

	expectedRbtcTotal := entities.NewWei(0)
	expectedRbtcTotal.Add(expectedRbtcTotal, expectedRbtcInRskWallet)
	expectedRbtcTotal.Add(expectedRbtcTotal, rbtcLockedInLbc)
	expectedRbtcTotal.Add(expectedRbtcTotal, expectedRbtcWaitingForRefund)

	expectedRbtcAvailable := entities.NewWei(0).Add(
		entities.NewWei(0).Sub(expectedRbtcInRskWallet, expectedRbtcReservedForUsers),
		rbtcLockedInLbc,
	)

	return ExpectedRbtcCalculations{
		RskWalletBalance: rbtcWalletBalance,
		InRskWallet:      expectedRbtcInRskWallet,
		LockedInLbc:      rbtcLockedInLbc,
		WaitingForRefund: expectedRbtcWaitingForRefund,
		ReservedForUsers: expectedRbtcReservedForUsers,
		Total:            expectedRbtcTotal,
		Available:        expectedRbtcAvailable,
	}
}

// Tests the BTC Asset Report generation with multiple pegout quote states to verify proper calculation of:
// - BTC Location: BtcWallet, Federation (rebalancing), RskWallet (waiting for rebalancing), Lbc
// - BTC Allocation: ReservedForUsers, WaitingForRefund, Available
// - Mathematical integrity of the BtcAssetReport
// nolint:exhaustive,funlen,maintidx
func TestGetAssetsReportUseCase_Run_BtcAssetReport_Success(t *testing.T) {
	testCases := []struct {
		name             string
		btcWalletBalance *entities.Wei
		pegoutQuotes     []quote.RetainedPegoutQuote
		description      string
	}{
		{
			name:             "Multiple quotes in various states",
			btcWalletBalance: entities.NewWei(50000000), // 0.5 BTC
			pegoutQuotes: []quote.RetainedPegoutQuote{
				// WaitingForDeposit quotes - should be counted in ReservedForUsers
				{QuoteHash: "waiting_deposit_1", RequiredLiquidity: entities.NewWei(1000000), State: quote.PegoutStateWaitingForDeposit},
				{QuoteHash: "waiting_deposit_2", RequiredLiquidity: entities.NewWei(2000000), State: quote.PegoutStateWaitingForDeposit},
				// WaitingForDepositConfirmations quotes - should be counted in ReservedForUsers
				{QuoteHash: "waiting_confirmations_1", RequiredLiquidity: entities.NewWei(1500000), State: quote.PegoutStateWaitingForDepositConfirmations},
				// BridgeTxSucceeded quotes - should be counted in Federation (rebalancing)
				{QuoteHash: "bridge_succeeded_1", RequiredLiquidity: entities.NewWei(5000000), State: quote.PegoutStateBridgeTxSucceeded},
				{QuoteHash: "bridge_succeeded_2", RequiredLiquidity: entities.NewWei(3000000), State: quote.PegoutStateBridgeTxSucceeded},
				// RefundPegOutSucceeded quotes - should be counted in RskWallet (waiting for rebalancing)
				{QuoteHash: "refund_succeeded_1", RequiredLiquidity: entities.NewWei(4000000), State: quote.PegoutStateRefundPegOutSucceeded},
				{QuoteHash: "refund_succeeded_2", RequiredLiquidity: entities.NewWei(2500000), State: quote.PegoutStateRefundPegOutSucceeded},
				// SendPegoutSucceeded quotes - should be counted in LBC (LP sent BTC, waiting for RBTC)
				{QuoteHash: "send_succeeded_1", RequiredLiquidity: entities.NewWei(3500000), State: quote.PegoutStateSendPegoutSucceeded},
				{QuoteHash: "send_succeeded_2", RequiredLiquidity: entities.NewWei(1800000), State: quote.PegoutStateSendPegoutSucceeded},
				// Other states that should not affect calculations
				{QuoteHash: "time_elapsed_1", RequiredLiquidity: entities.NewWei(1000000), State: quote.PegoutStateTimeForDepositElapsed},
				{QuoteHash: "send_failed_1", RequiredLiquidity: entities.NewWei(2000000), State: quote.PegoutStateSendPegoutFailed},
			},
			description: "Tests with quotes in all relevant states including waiting, rebalancing, and completed states",
		},
		{
			name:             "Only waiting for deposit quotes",
			btcWalletBalance: entities.NewWei(100000000), // 1.0 BTC
			pegoutQuotes: []quote.RetainedPegoutQuote{
				{
					QuoteHash:         "waiting_1",
					RequiredLiquidity: entities.NewWei(5000000),
					State:             quote.PegoutStateWaitingForDeposit,
				},
				{
					QuoteHash:         "waiting_2",
					RequiredLiquidity: entities.NewWei(3000000),
					State:             quote.PegoutStateWaitingForDeposit,
				},
				{
					QuoteHash:         "waiting_confirmations_1",
					RequiredLiquidity: entities.NewWei(2000000),
					State:             quote.PegoutStateWaitingForDepositConfirmations,
				},
			},
			description: "All BTC should be in wallet with reserved amount, large available balance",
		},
		{
			name:             "Only rebalancing quotes",
			btcWalletBalance: entities.NewWei(30000000), // 0.3 BTC
			pegoutQuotes: []quote.RetainedPegoutQuote{
				{
					QuoteHash:         "bridge_1",
					RequiredLiquidity: entities.NewWei(10000000),
					State:             quote.PegoutStateBridgeTxSucceeded,
				},
				{
					QuoteHash:         "bridge_2",
					RequiredLiquidity: entities.NewWei(15000000),
					State:             quote.PegoutStateBridgeTxSucceeded,
				},
			},
			description: "BTC distributed between wallet and federation (rebalancing)",
		},
		{
			name:             "Mix of refund and send states",
			btcWalletBalance: entities.NewWei(20000000), // 0.2 BTC
			pegoutQuotes: []quote.RetainedPegoutQuote{
				{
					QuoteHash:         "refund_1",
					RequiredLiquidity: entities.NewWei(8000000),
					State:             quote.PegoutStateRefundPegOutSucceeded,
				},
				{
					QuoteHash:         "send_1",
					RequiredLiquidity: entities.NewWei(5000000),
					State:             quote.PegoutStateSendPegoutSucceeded,
				},
				{
					QuoteHash:         "send_2",
					RequiredLiquidity: entities.NewWei(7000000),
					State:             quote.PegoutStateSendPegoutSucceeded,
				},
			},
			description: "BTC in wallet, RSK wallet (waiting for rebalancing), and LBC",
		},
		{
			name:             "Empty quotes - only wallet balance",
			btcWalletBalance: entities.NewWei(75000000), // 0.75 BTC
			pegoutQuotes:     []quote.RetainedPegoutQuote{},
			description:      "All BTC in wallet, no quotes, everything available",
		},
		{
			name:             "All states combined - complex scenario",
			btcWalletBalance: entities.NewWei(200000000), // 2.0 BTC
			pegoutQuotes: []quote.RetainedPegoutQuote{
				{QuoteHash: "waiting_1", RequiredLiquidity: entities.NewWei(10000000), State: quote.PegoutStateWaitingForDeposit},
				{QuoteHash: "waiting_2", RequiredLiquidity: entities.NewWei(5000000), State: quote.PegoutStateWaitingForDepositConfirmations},
				{QuoteHash: "bridge_1", RequiredLiquidity: entities.NewWei(25000000), State: quote.PegoutStateBridgeTxSucceeded},
				{QuoteHash: "bridge_2", RequiredLiquidity: entities.NewWei(30000000), State: quote.PegoutStateBridgeTxSucceeded},
				{QuoteHash: "refund_1", RequiredLiquidity: entities.NewWei(20000000), State: quote.PegoutStateRefundPegOutSucceeded},
				{QuoteHash: "refund_2", RequiredLiquidity: entities.NewWei(15000000), State: quote.PegoutStateRefundPegOutSucceeded},
				{QuoteHash: "send_1", RequiredLiquidity: entities.NewWei(18000000), State: quote.PegoutStateSendPegoutSucceeded},
				{QuoteHash: "send_2", RequiredLiquidity: entities.NewWei(12000000), State: quote.PegoutStateSendPegoutSucceeded},
			},
			description: "Large wallet with quotes in all states - comprehensive test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			expected := calculateExpectedBtcValues(tc.pegoutQuotes, tc.btcWalletBalance)

			btcWallet := &mocks.BitcoinWalletMock{}
			rskRpc := &mocks.RootstockRpcServerMock{}
			lp := &mocks.ProviderMock{}
			peginProvider := &mocks.ProviderMock{}
			pegoutProvider := &mocks.ProviderMock{}
			peginRepository := &mocks.PeginQuoteRepositoryMock{}
			pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
			peginContract := &mocks.PeginContractMock{}

			contracts := blockchain.RskContracts{
				PegIn: peginContract,
			}

			btcWallet.On("GetBalance").Return(tc.btcWalletBalance, nil).Once()

			// Separate quotes by state for repository expectations
			bridgeQuotes := []quote.RetainedPegoutQuote{}
			refundQuotes := []quote.RetainedPegoutQuote{}
			sendQuotes := []quote.RetainedPegoutQuote{}
			waitingQuotes := []quote.RetainedPegoutQuote{}
			combinedWaitingForRefundQuotes := []quote.RetainedPegoutQuote{}

			for _, q := range tc.pegoutQuotes {
				switch q.State {
				case quote.PegoutStateBridgeTxSucceeded:
					bridgeQuotes = append(bridgeQuotes, q)
					combinedWaitingForRefundQuotes = append(combinedWaitingForRefundQuotes, q)
				case quote.PegoutStateRefundPegOutSucceeded:
					refundQuotes = append(refundQuotes, q)
					combinedWaitingForRefundQuotes = append(combinedWaitingForRefundQuotes, q)
				case quote.PegoutStateSendPegoutSucceeded:
					sendQuotes = append(sendQuotes, q)
					combinedWaitingForRefundQuotes = append(combinedWaitingForRefundQuotes, q)
				case quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations:
					waitingQuotes = append(waitingQuotes, q)
				}
			}

			// BridgeTxSucceeded state (rebalancing)
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
				Return(retainedPegoutQuotesToPegoutQuotes(t, bridgeQuotes), nil).Once()

			// RefundPegOutSucceeded state (waiting for rebalancing)
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
				Return(retainedPegoutQuotesToPegoutQuotes(t, refundQuotes), nil).Once()

			// SendPegoutSucceeded state (in LBC)
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
				Return(retainedPegoutQuotesToPegoutQuotes(t, sendQuotes), nil).Once()

			// WaitingForDeposit and WaitingForDepositConfirmations states (reserved for users)
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
				Return(retainedPegoutQuotesToPegoutQuotes(t, waitingQuotes), nil).Once()

			// Combined states for waiting for refund calculation
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
				Return(retainedPegoutQuotesToPegoutQuotes(t, combinedWaitingForRefundQuotes), nil).Once()

			// Setup mock expectations for RBTC-related calls (minimal setup to make them pass)
			rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(entities.NewWei(0), nil).Once()
			lp.On("RskAddress").Return("test-rsk-address").Twice() // Called twice: once for RSK balance, once for LBC balance
			peginContract.On("GetBalance", "test-rsk-address").Return(entities.NewWei(0), nil).Once()
			peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
				Return([]quote.PeginQuote{}, nil).Once()
			peginRepository.On("GetQuotesByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).
				Return([]quote.PeginQuote{}, nil).Once()

			useCase := reports.NewGetAssetsReportUseCase(
				btcWallet,
				blockchain.Rpc{Rsk: rskRpc},
				lp,
				peginProvider,
				pegoutProvider,
				peginRepository,
				pegoutRepository,
				contracts,
			)

			result, err := useCase.Run(ctx)

			require.NoError(t, err)

			assert.Equal(t, expected.Total.String(), result.BtcAssetReport.Total.String(), "BTC total should match expected")
			assert.Equal(t, expected.WalletBalance.String(), result.BtcAssetReport.Location.BtcWallet.String(), "BTC wallet balance should match")
			assert.Equal(t, expected.Rebalancing.String(), result.BtcAssetReport.Location.Federation.String(), "BTC federation balance should match")
			assert.Equal(t, expected.WaitingForRebalancing.String(), result.BtcAssetReport.Location.RskWallet.String(), "BTC RSK wallet balance should match")
			assert.Equal(t, expected.InLbc.String(), result.BtcAssetReport.Location.Lbc.String(), "BTC LBC balance should match")
			assert.Equal(t, expected.ReservedForUsers.String(), result.BtcAssetReport.Allocation.ReservedForUsers.String(), "BTC reserved for users should match")
			assert.Equal(t, expected.WaitingForRefund.String(), result.BtcAssetReport.Allocation.WaitingForRefund.String(), "BTC waiting for refund should match")
			assert.Equal(t, expected.Available.String(), result.BtcAssetReport.Allocation.Available.String(), "BTC available should match")

			locationSum := entities.NewWei(0)
			locationSum.Add(locationSum, result.BtcAssetReport.Location.BtcWallet)
			locationSum.Add(locationSum, result.BtcAssetReport.Location.Federation)
			locationSum.Add(locationSum, result.BtcAssetReport.Location.RskWallet)
			locationSum.Add(locationSum, result.BtcAssetReport.Location.Lbc)
			assert.Equal(t, result.BtcAssetReport.Total.String(), locationSum.String(), "Location sum should equal Total")

			allocationSum := entities.NewWei(0)
			allocationSum.Add(allocationSum, result.BtcAssetReport.Allocation.ReservedForUsers)
			allocationSum.Add(allocationSum, result.BtcAssetReport.Allocation.WaitingForRefund)
			allocationSum.Add(allocationSum, result.BtcAssetReport.Allocation.Available)
			assert.Equal(t, result.BtcAssetReport.Total.String(), allocationSum.String(), "Allocation sum should equal Total")

			btcWallet.AssertExpectations(t)
			rskRpc.AssertExpectations(t)
			lp.AssertExpectations(t)
			peginRepository.AssertExpectations(t)
			pegoutRepository.AssertExpectations(t)
			peginContract.AssertExpectations(t)
		})
	}
}

// Tests the RBTC Asset Report generation with multiple pegin quote states to verify proper calculation of:
// - RBTC Location: RskWallet, Lbc, Federation (waiting for refund)
// - RBTC Allocation: ReservedForUsers, WaitingForRefund, Available
// - Mathematical integrity of the RbtcAssetReport
// nolint:funlen,exhaustive
func TestGetAssetsReportUseCase_Run_RbtcAssetReport_Success(t *testing.T) {
	testCases := []struct {
		name                     string
		rbtcWalletBalance        *entities.Wei
		rbtcLockedInLbc          *entities.Wei
		peginQuotes              []quote.RetainedPeginQuote
		btcWaitingForRebalancing *entities.Wei // BTC value needed for RBTC calculation
		description              string
	}{
		{
			name:                     "Multiple pegin quotes in various states",
			rbtcWalletBalance:        entities.NewWei(50000000), // 0.5 RBTC in RSK wallet
			rbtcLockedInLbc:          entities.NewWei(30000000), // 0.3 RBTC in LBC
			btcWaitingForRebalancing: entities.NewWei(10000000), // 0.1 BTC waiting for rebalancing
			peginQuotes: []quote.RetainedPeginQuote{
				// CallForUserSucceeded quotes - should be counted in WaitingForRefund (Federation)
				{QuoteHash: "call_succeeded_1", RequiredLiquidity: entities.NewWei(5000000), State: quote.PeginStateCallForUserSucceeded},
				{QuoteHash: "call_succeeded_2", RequiredLiquidity: entities.NewWei(3000000), State: quote.PeginStateCallForUserSucceeded},
				// WaitingForDeposit quotes - should be counted in ReservedForUsers
				{QuoteHash: "pegin_waiting_1", RequiredLiquidity: entities.NewWei(7000000), State: quote.PeginStateWaitingForDeposit},
				// WaitingForDepositConfirmations quotes - should be counted in ReservedForUsers
				{QuoteHash: "pegin_waiting_conf_1", RequiredLiquidity: entities.NewWei(4000000), State: quote.PeginStateWaitingForDepositConfirmations},
				// Other states that should not affect calculations
				{QuoteHash: "time_elapsed_1", RequiredLiquidity: entities.NewWei(2000000), State: quote.PeginStateTimeForDepositElapsed},
				{QuoteHash: "register_failed_1", RequiredLiquidity: entities.NewWei(1000000), State: quote.PeginStateRegisterPegInFailed},
			},
			description: "Tests with pegin quotes in all relevant states including waiting and completed states",
		},
		{
			name:                     "Only waiting for deposit quotes",
			rbtcWalletBalance:        entities.NewWei(100000000), // 1.0 RBTC
			rbtcLockedInLbc:          entities.NewWei(50000000),  // 0.5 RBTC
			btcWaitingForRebalancing: entities.NewWei(0),         // No BTC waiting
			peginQuotes: []quote.RetainedPeginQuote{
				{QuoteHash: "waiting_1", RequiredLiquidity: entities.NewWei(15000000), State: quote.PeginStateWaitingForDeposit},
				{QuoteHash: "waiting_2", RequiredLiquidity: entities.NewWei(10000000), State: quote.PeginStateWaitingForDeposit},
				{QuoteHash: "waiting_conf_1", RequiredLiquidity: entities.NewWei(5000000), State: quote.PeginStateWaitingForDepositConfirmations},
			},
			description: "All RBTC in wallet and LBC with reserved amount, large available balance",
		},
		{
			name:                     "Only CallForUserSucceeded quotes",
			rbtcWalletBalance:        entities.NewWei(80000000), // 0.8 RBTC
			rbtcLockedInLbc:          entities.NewWei(40000000), // 0.4 RBTC
			btcWaitingForRebalancing: entities.NewWei(5000000),  // 0.05 BTC waiting
			peginQuotes: []quote.RetainedPeginQuote{
				{QuoteHash: "call_1", RequiredLiquidity: entities.NewWei(12000000), State: quote.PeginStateCallForUserSucceeded},
				{QuoteHash: "call_2", RequiredLiquidity: entities.NewWei(8000000), State: quote.PeginStateCallForUserSucceeded},
			},
			description: "RBTC waiting for refund in federation",
		},
		{
			name:                     "Empty quotes - only wallet and LBC balance",
			rbtcWalletBalance:        entities.NewWei(120000000), // 1.2 RBTC
			rbtcLockedInLbc:          entities.NewWei(80000000),  // 0.8 RBTC
			btcWaitingForRebalancing: entities.NewWei(20000000),  // 0.2 BTC waiting
			peginQuotes:              []quote.RetainedPeginQuote{},
			description:              "All RBTC available, no quotes",
		},
		{
			name:                     "Large BTC waiting for rebalancing affects RSK wallet",
			rbtcWalletBalance:        entities.NewWei(150000000), // 1.5 RBTC raw balance
			rbtcLockedInLbc:          entities.NewWei(60000000),  // 0.6 RBTC
			btcWaitingForRebalancing: entities.NewWei(50000000),  // 0.5 BTC waiting (reduces effective RSK wallet)
			peginQuotes: []quote.RetainedPeginQuote{
				{QuoteHash: "waiting_1", RequiredLiquidity: entities.NewWei(20000000), State: quote.PeginStateWaitingForDeposit},
				{QuoteHash: "call_1", RequiredLiquidity: entities.NewWei(15000000), State: quote.PeginStateCallForUserSucceeded},
			},
			description: "BTC waiting for rebalancing reduces effective RBTC in RSK wallet",
		},
		{
			name:                     "All states combined - complex scenario",
			rbtcWalletBalance:        entities.NewWei(200000000), // 2.0 RBTC
			rbtcLockedInLbc:          entities.NewWei(100000000), // 1.0 RBTC
			btcWaitingForRebalancing: entities.NewWei(30000000),  // 0.3 BTC waiting
			peginQuotes: []quote.RetainedPeginQuote{
				{QuoteHash: "waiting_1", RequiredLiquidity: entities.NewWei(25000000), State: quote.PeginStateWaitingForDeposit},
				{QuoteHash: "waiting_2", RequiredLiquidity: entities.NewWei(15000000), State: quote.PeginStateWaitingForDepositConfirmations},
				{QuoteHash: "call_1", RequiredLiquidity: entities.NewWei(35000000), State: quote.PeginStateCallForUserSucceeded},
				{QuoteHash: "call_2", RequiredLiquidity: entities.NewWei(20000000), State: quote.PeginStateCallForUserSucceeded},
				{QuoteHash: "call_3", RequiredLiquidity: entities.NewWei(10000000), State: quote.PeginStateCallForUserSucceeded},
			},
			description: "Large wallet with pegin quotes in all states - comprehensive test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			expectedRbtc := calculateExpectedRbtcValues(tc.peginQuotes, tc.rbtcWalletBalance, tc.rbtcLockedInLbc, tc.btcWaitingForRebalancing)

			btcWallet := &mocks.BitcoinWalletMock{}
			rskRpc := &mocks.RootstockRpcServerMock{}
			lp := &mocks.ProviderMock{}
			peginProvider := &mocks.ProviderMock{}
			pegoutProvider := &mocks.ProviderMock{}
			peginRepository := &mocks.PeginQuoteRepositoryMock{}
			pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
			peginContract := &mocks.PeginContractMock{}
			contracts := blockchain.RskContracts{
				PegIn: peginContract,
			}

			// Setup mock expectations for BTC-related calls (minimal setup to make them pass)
			// Note: We need to mock pegout quotes to calculate btcWaitingForRebalancing
			btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()

			// Create pegout quotes that will result in the desired btcWaitingForRebalancing
			// Note: The quote's Total() = RequiredLiquidity + additionalFees, so we subtract the fees
			// to ensure the final Total() equals the desired btcWaitingForRebalancing amount
			pegoutQuotesForRebalancing := []quote.RetainedPegoutQuote{}
			if tc.btcWaitingForRebalancing.Cmp(entities.NewWei(0)) > 0 {
				// Subtract additional fees from the required liquidity so that Total() equals btcWaitingForRebalancing
				requiredLiquidityForQuote := new(entities.Wei).Sub(tc.btcWaitingForRebalancing, additionalFeesPerQuote())
				pegoutQuotesForRebalancing = append(pegoutQuotesForRebalancing, quote.RetainedPegoutQuote{
					QuoteHash:         "refund_for_test",
					RequiredLiquidity: requiredLiquidityForQuote,
					State:             quote.PegoutStateRefundPegOutSucceeded,
				})
			}

			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
				Return([]quote.PegoutQuote{}, nil).Once()
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
				Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesForRebalancing), nil).Once()
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
				Return([]quote.PegoutQuote{}, nil).Once()
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
				Return([]quote.PegoutQuote{}, nil).Once()
			pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
				Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesForRebalancing), nil).Once()

			// Setup mock expectations for RBTC-related calls
			rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(tc.rbtcWalletBalance, nil).Once()
			lp.On("RskAddress").Return("test-rsk-address").Twice() // Called twice: once for RSK balance, once for LBC balance
			peginContract.On("GetBalance", "test-rsk-address").Return(tc.rbtcLockedInLbc, nil).Once()

			// Setup mock expectations for pegin quotes by different states
			callSucceededQuotes := []quote.RetainedPeginQuote{}
			waitingQuotes := []quote.RetainedPeginQuote{}

			for _, q := range tc.peginQuotes {
				switch q.State {
				case quote.PeginStateCallForUserSucceeded:
					callSucceededQuotes = append(callSucceededQuotes, q)
				case quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations:
					waitingQuotes = append(waitingQuotes, q)
				}
			}

			peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
				Return(retainedPeginQuotesToPeginQuotes(t, callSucceededQuotes), nil).Once()
			peginRepository.On("GetQuotesByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).
				Return(retainedPeginQuotesToPeginQuotes(t, waitingQuotes), nil).Once()

			useCase := reports.NewGetAssetsReportUseCase(
				btcWallet,
				blockchain.Rpc{Rsk: rskRpc},
				lp,
				peginProvider,
				pegoutProvider,
				peginRepository,
				pegoutRepository,
				contracts,
			)

			result, err := useCase.Run(ctx)

			require.NoError(t, err)

			assert.Equal(t, expectedRbtc.Total.String(), result.RbtcAssetReport.Total.String(), "RBTC total should match expected")
			assert.Equal(t, expectedRbtc.InRskWallet.String(), result.RbtcAssetReport.Location.RskWallet.String(), "RBTC RSK wallet balance should match")
			assert.Equal(t, expectedRbtc.LockedInLbc.String(), result.RbtcAssetReport.Location.Lbc.String(), "RBTC LBC balance should match")
			assert.Equal(t, expectedRbtc.WaitingForRefund.String(), result.RbtcAssetReport.Location.Federation.String(), "RBTC federation balance should match")
			assert.Equal(t, expectedRbtc.ReservedForUsers.String(), result.RbtcAssetReport.Allocation.ReservedForUsers.String(), "RBTC reserved for users should match")
			assert.Equal(t, expectedRbtc.WaitingForRefund.String(), result.RbtcAssetReport.Allocation.WaitingForRefund.String(), "RBTC waiting for refund should match")
			assert.Equal(t, expectedRbtc.Available.String(), result.RbtcAssetReport.Allocation.Available.String(), "RBTC available should match")

			locationSum := entities.NewWei(0)
			locationSum.Add(locationSum, result.RbtcAssetReport.Location.RskWallet)
			locationSum.Add(locationSum, result.RbtcAssetReport.Location.Lbc)
			locationSum.Add(locationSum, result.RbtcAssetReport.Location.Federation)
			assert.Equal(t, result.RbtcAssetReport.Total.String(), locationSum.String(), "RBTC Location sum should equal Total")

			allocationSum := entities.NewWei(0)
			allocationSum.Add(allocationSum, result.RbtcAssetReport.Allocation.ReservedForUsers)
			allocationSum.Add(allocationSum, result.RbtcAssetReport.Allocation.WaitingForRefund)
			allocationSum.Add(allocationSum, result.RbtcAssetReport.Allocation.Available)
			assert.Equal(t, result.RbtcAssetReport.Total.String(), allocationSum.String(), "RBTC Allocation sum should equal Total")

			btcWallet.AssertExpectations(t)
			rskRpc.AssertExpectations(t)
			lp.AssertExpectations(t)
			peginRepository.AssertExpectations(t)
			pegoutRepository.AssertExpectations(t)
			peginContract.AssertExpectations(t)
		})
	}
}

// Tests that total assets remain constant as quotes progress through their lifecycle states.
// This verifies that the report correctly tracks assets regardless of quote state changes.
// nolint:funlen
func TestGetAssetsReportUseCase_Run_AssetConservation_ThroughQuoteLifecycle(t *testing.T) {
	ctx := context.Background()

	initialBtcWalletBalance := entities.NewWei(100000000)  // 1.0 BTC
	initialRbtcWalletBalance := entities.NewWei(200000000) // 2.0 RBTC
	initialRbtcLbcBalance := entities.NewWei(50000000)     // 0.5 RBTC

	// Calculate initial total assets (without fees for simplicity)
	expectedTotalBtc := entities.NewWei(0).Add(entities.NewWei(0), initialBtcWalletBalance)
	expectedTotalRbtc := entities.NewWei(0).Add(entities.NewWei(0), initialRbtcWalletBalance)
	expectedTotalRbtc.Add(expectedTotalRbtc, initialRbtcLbcBalance)

	pegoutQuote1Amount := entities.NewWei(10000000) // 0.1 BTC
	pegoutQuote2Amount := entities.NewWei(15000000) // 0.15 BTC
	peginQuote1Amount := entities.NewWei(20000000)  // 0.2 RBTC
	peginQuote2Amount := entities.NewWei(25000000)  // 0.25 RBTC

	// Scenario 1: Initial state - quotes just accepted (waiting for deposit)
	t.Run("Scenario_1_Quotes_Waiting_For_Deposit", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		pegoutQuotes := []quote.RetainedPegoutQuote{
			{QuoteHash: "pegout_1", RequiredLiquidity: pegoutQuote1Amount, State: quote.PegoutStateWaitingForDeposit},
			{QuoteHash: "pegout_2", RequiredLiquidity: pegoutQuote2Amount, State: quote.PegoutStateWaitingForDepositConfirmations},
		}
		peginQuotes := []quote.RetainedPeginQuote{
			{QuoteHash: "pegin_1", RequiredLiquidity: peginQuote1Amount, State: quote.PeginStateWaitingForDeposit},
			{QuoteHash: "pegin_2", RequiredLiquidity: peginQuote2Amount, State: quote.PeginStateWaitingForDepositConfirmations},
		}

		btcWallet.On("GetBalance").Return(initialBtcWalletBalance, nil).Once()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(initialRbtcWalletBalance, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Twice()
		peginContract.On("GetBalance", "test-rsk-address").Return(initialRbtcLbcBalance, nil).Once()

		// Pegout repository mocks
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotes), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()

		// Pegin repository mocks
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
			Return([]quote.PeginQuote{}, nil).Once()
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).
			Return(retainedPeginQuotesToPeginQuotes(t, peginQuotes), nil).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedTotalBtc.String(), result.BtcAssetReport.Total.String(), "Scenario 1: BTC total should remain constant")
		assert.Equal(t, expectedTotalRbtc.String(), result.RbtcAssetReport.Total.String(), "Scenario 1: RBTC total should remain constant")

		btcWallet.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
		peginRepository.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		peginContract.AssertExpectations(t)
	})

	// Scenario 2: Quotes progress - LP sends BTC (pegout), LP calls for user (pegin)
	t.Run("Scenario_2_Quotes_In_Progress", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		// Scenario 2: User deposits completed, assets moved but not yet fully cycled
		// - pegout_1: User deposited RBTC, LP refunded (RefundPegOutSucceeded) - LP now has RBTC in wallet
		//   The report counts this RBTC as "BTC waiting for rebalancing" because LP will convert it back to BTC
		// - pegout_2: Still waiting for user deposit
		// - pegin_1: LP called for user (CallForUserSucceeded) - LP sent RBTC, waiting for BTC from user
		// - pegin_2: Still waiting for user deposit

		pegoutQuotesRefunded := []quote.RetainedPegoutQuote{
			{QuoteHash: "pegout_1", RequiredLiquidity: pegoutQuote1Amount, State: quote.PegoutStateRefundPegOutSucceeded},
		}
		pegoutQuotesReserved := []quote.RetainedPegoutQuote{
			{QuoteHash: "pegout_2", RequiredLiquidity: pegoutQuote2Amount, State: quote.PegoutStateWaitingForDepositConfirmations},
		}
		peginQuotesWaitingRefund := []quote.RetainedPeginQuote{
			{QuoteHash: "pegin_1", RequiredLiquidity: peginQuote1Amount, State: quote.PeginStateCallForUserSucceeded},
		}
		peginQuotesReserved := []quote.RetainedPeginQuote{
			{QuoteHash: "pegin_2", RequiredLiquidity: peginQuote2Amount, State: quote.PeginStateWaitingForDepositConfirmations},
		}

		// RBTC wallet balance changes:
		// - Reduced by pegin_1 amount (sent to user in callForUser)
		// - Increased by pegout_1 amount (received from user's RBTC deposit)
		// Net: initialRbtcWalletBalance - peginQuote1Amount + pegoutQuote1Amount
		currentRbtcWalletBalance := entities.NewWei(0).Add(initialRbtcWalletBalance, entities.NewWei(0))
		currentRbtcWalletBalance.Sub(currentRbtcWalletBalance, peginQuote1Amount)
		currentRbtcWalletBalance.Add(currentRbtcWalletBalance, pegoutQuote1Amount)

		// Expected totals adjust for the asset type conversion:
		// BTC total increases by pegoutQuote1Amount + fees (RBTC counted as "BTC waiting for rebalancing")
		// RBTC total stays the same (pegin_1 sent out, pegout_1 received in)
		pegoutQuote1Total := new(entities.Wei).Add(pegoutQuote1Amount, additionalFeesPerQuote())
		expectedBtcScenario2 := entities.NewWei(0).Add(expectedTotalBtc, pegoutQuote1Total)

		btcWallet.On("GetBalance").Return(initialBtcWalletBalance, nil).Once()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(currentRbtcWalletBalance, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Twice()
		peginContract.On("GetBalance", "test-rsk-address").Return(initialRbtcLbcBalance, nil).Once()

		// Pegout repository mocks
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesRefunded), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesReserved), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesRefunded), nil).Once()

		// Pegin repository mocks
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
			Return(retainedPeginQuotesToPeginQuotes(t, peginQuotesWaitingRefund), nil).Once()
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).
			Return(retainedPeginQuotesToPeginQuotes(t, peginQuotesReserved), nil).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedBtcScenario2.String(), result.BtcAssetReport.Total.String(), "Scenario 2: BTC total includes RBTC waiting for rebalancing")
		assert.Equal(t, expectedTotalRbtc.String(), result.RbtcAssetReport.Total.String(), "Scenario 2: RBTC total remains constant")

		btcWallet.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
		peginRepository.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		peginContract.AssertExpectations(t)
	})

	// Scenario 3: Final state - quotes completed and refunded
	t.Run("Scenario_3_Quotes_Completed_And_Refunded", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		// Scenario 3: All quotes completed and refunded
		// - pegout_1: LP sent BTC, user deposited RBTC to LBC, LP refunded (RefundPegOutSucceeded)
		// - pegout_2: Still waiting for user deposit (WaitingForDepositConfirmations)
		// - pegin_1: LP called for user, user sent BTC, LP registered and got refunded - quote completed (no longer in system)
		// - pegin_2: Still waiting for user deposit (WaitingForDepositConfirmations)

		pegoutQuotesRefunded := []quote.RetainedPegoutQuote{
			{QuoteHash: "pegout_1", RequiredLiquidity: pegoutQuote1Amount, State: quote.PegoutStateRefundPegOutSucceeded},
		}
		pegoutQuotesStillWaiting := []quote.RetainedPegoutQuote{
			{QuoteHash: "pegout_2", RequiredLiquidity: pegoutQuote2Amount, State: quote.PegoutStateWaitingForDepositConfirmations},
		}
		peginQuotesStillWaiting := []quote.RetainedPeginQuote{
			{QuoteHash: "pegin_2", RequiredLiquidity: peginQuote2Amount, State: quote.PeginStateWaitingForDepositConfirmations},
		}

		// LP got RBTC back from pegin_1 refund, now has the RBTC in RSK wallet (not in LBC anymore, assuming LP withdrew)
		// The pegin_1 amount is back in the wallet
		finalRbtcWalletBalance := entities.NewWei(0).Add(initialRbtcWalletBalance, entities.NewWei(0))
		// LP has the RBTC from pegout_1 in RSK wallet (received from refund)
		finalRbtcWalletBalance.Add(finalRbtcWalletBalance, pegoutQuote1Amount)

		btcWallet.On("GetBalance").Return(initialBtcWalletBalance, nil).Once()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(finalRbtcWalletBalance, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Twice()
		peginContract.On("GetBalance", "test-rsk-address").Return(initialRbtcLbcBalance, nil).Once()

		// Pegout repository mocks
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesRefunded), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesStillWaiting), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return(retainedPegoutQuotesToPegoutQuotes(t, pegoutQuotesRefunded), nil).Once()

		// Pegin repository mocks
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
			Return([]quote.PeginQuote{}, nil).Once()
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).
			Return(retainedPeginQuotesToPeginQuotes(t, peginQuotesStillWaiting), nil).Once()

		// Expected totals: Same as Scenario 2 since pegout_1 is still in RefundPegOutSucceeded state
		// Include additional fees in the quote total
		pegoutQuote1Total := new(entities.Wei).Add(pegoutQuote1Amount, additionalFeesPerQuote())
		expectedBtcScenario3 := entities.NewWei(0).Add(expectedTotalBtc, pegoutQuote1Total)

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedBtcScenario3.String(), result.BtcAssetReport.Total.String(), "Scenario 3: BTC total includes RBTC waiting for rebalancing")
		// In Scenario 3, pegin_1 is complete (no longer waiting), so RBTC total is reduced by the fees that were part of pegin_1's Total()
		expectedRbtcScenario3 := new(entities.Wei).Sub(expectedTotalRbtc, additionalFeesPerQuote())
		assert.Equal(t, expectedRbtcScenario3.String(), result.RbtcAssetReport.Total.String(), "Scenario 3: RBTC total adjusted for completed pegin_1")

		btcWallet.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
		peginRepository.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		peginContract.AssertExpectations(t)
	})
}

// Tests error handling in the GetAssetsReportUseCase for various failure scenarios.
// Verifies that errors from dependencies are properly propagated.
// nolint:funlen,maintidx
func TestGetAssetsReportUseCase_Run_ErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("Error_BtcWallet_GetBalance_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
	})

	t.Run("Error_PegoutRepository_BridgeTxSucceeded_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
	})

	t.Run("Error_PegoutRepository_RefundPegOutSucceeded_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
	})

	t.Run("Error_PegoutRepository_SendPegoutSucceeded_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
	})

	t.Run("Error_PegoutRepository_WaitingForDeposit_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
	})

	t.Run("Error_PegoutRepository_WaitingForRefund_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
	})

	t.Run("Error_RskRpc_GetBalance_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Once()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
	})

	t.Run("Error_LbcContract_GetBalance_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Twice()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(entities.NewWei(200000000), nil).Once()
		peginContract.On("GetBalance", "test-rsk-address").Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
		peginContract.AssertExpectations(t)
	})

	t.Run("Error_PeginRepository_CallForUserSucceeded_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Twice()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(entities.NewWei(200000000), nil).Once()
		peginContract.On("GetBalance", "test-rsk-address").Return(entities.NewWei(50000000), nil).Once()
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
		peginContract.AssertExpectations(t)
		peginRepository.AssertExpectations(t)
	})

	t.Run("Error_PeginRepository_WaitingForDeposit_Fails", func(t *testing.T) {
		btcWallet := &mocks.BitcoinWalletMock{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		lp := &mocks.ProviderMock{}
		peginProvider := &mocks.ProviderMock{}
		pegoutProvider := &mocks.ProviderMock{}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		peginContract := &mocks.PeginContractMock{}
		contracts := blockchain.RskContracts{PegIn: peginContract}

		btcWallet.On("GetBalance").Return(entities.NewWei(100000000), nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateSendPegoutSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.PegoutQuote{}, nil).Once()
		pegoutRepository.On("GetQuotesByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return([]quote.PegoutQuote{}, nil).Once()
		lp.On("RskAddress").Return("test-rsk-address").Twice()
		rskRpc.On("GetBalance", ctx, "test-rsk-address").Return(entities.NewWei(200000000), nil).Once()
		peginContract.On("GetBalance", "test-rsk-address").Return(entities.NewWei(50000000), nil).Once()
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateCallForUserSucceeded).
			Return([]quote.PeginQuote{}, nil).Once()
		peginRepository.On("GetQuotesByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetAssetsReportUseCase(btcWallet, blockchain.Rpc{Rsk: rskRpc}, lp, peginProvider, pegoutProvider, peginRepository, pegoutRepository, contracts)
		result, err := useCase.Run(ctx)

		require.Error(t, err)
		assert.Equal(t, reports.GetAssetsReportResult{}, result)
		btcWallet.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
		lp.AssertExpectations(t)
		peginContract.AssertExpectations(t)
		peginRepository.AssertExpectations(t)
	})
}
