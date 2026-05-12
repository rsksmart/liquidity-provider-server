package contracts

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/integration"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SplitLbcExecutor struct {
	pegin                *bindings.IPegIn
	pegout               *bindings.IPegOut
	discovery            *bindings.IFlyoverDiscovery
	collateralManagement *bindings.ICollateralManagement
}

type SplitAddresses struct {
	Discovery            string
	Pegout               string
	Pegin                string
	CollateralManagement string
}

func NewSplityLbcExecutor(addresses SplitAddresses, backend bind.ContractBackend) (*SplitLbcExecutor, error) {
	var (
		err                  error
		pegin                *bindings.IPegIn
		pegout               *bindings.IPegOut
		discovery            *bindings.IFlyoverDiscovery
		collateralManagement *bindings.ICollateralManagement
	)
	if pegin, err = bindings.NewIPegIn(common.HexToAddress(addresses.Pegin), backend); err != nil {
		return nil, err
	}
	if pegout, err = bindings.NewIPegOut(common.HexToAddress(addresses.Pegout), backend); err != nil {
		return nil, err
	}
	if discovery, err = bindings.NewIFlyoverDiscovery(common.HexToAddress(addresses.Discovery), backend); err != nil {
		return nil, err
	}
	if collateralManagement, err = bindings.NewICollateralManagement(common.HexToAddress(addresses.CollateralManagement), backend); err != nil {
		return nil, err
	}
	return &SplitLbcExecutor{pegin: pegin, pegout: pegout, discovery: discovery, collateralManagement: collateralManagement}, nil
}

func (e *SplitLbcExecutor) DepositPegout(s TestSuite, opts *bind.TransactOpts, pegoutQuote pkg.PegoutQuoteDTO, hexSignature string) (*types.Receipt, *types.Transaction) {
	opts.Value = integration.SumAll(
		pegoutQuote.Value,
		pegoutQuote.CallFee,
		pegoutQuote.GasFee,
	)

	ctx := context.Background()
	gasPrice, err := s.RskClient().SuggestGasPrice(ctx)
	s.Raw().Require().NoError(err)
	opts.GasPrice = gasPrice

	parsedQuote := e.parsePegoutQuote(s.Raw(), pegoutQuote)
	signature, err := hex.DecodeString(hexSignature)
	s.Raw().Require().NoError(err)

	depositTx, err := e.pegout.DepositPegOut(opts, parsedQuote, signature)
	s.Raw().Require().NoError(err)
	receipt, err := bind.WaitMined(ctx, s.RskClient(), depositTx)
	s.Raw().Require().NoError(err)
	log.Info("[Integration test] Hash of deposit tx ", depositTx.Hash().String())
	return receipt, depositTx
}

func (e *SplitLbcExecutor) GetRefundPegoutEvent(s TestSuite, timeout time.Duration, quoteHash string) RefundPegoutEvent {
	var quoteHashByes [32]byte
	eventChannel := make(chan *bindings.IPegOutPegOutRefunded)
	parsedHash, err := hex.DecodeString(quoteHash)
	s.Raw().Require().NoError(err)
	copy(quoteHashByes[:], parsedHash)

	subscription, err := e.pegout.WatchPegOutRefunded(
		nil,
		eventChannel,
		[][32]byte{quoteHashByes},
	)
	s.Raw().Require().NoError(err)
	defer subscription.Unsubscribe()

	done := make(chan os.Signal, 1)
	testTolerance := time.NewTimer(timeout)
	defer testTolerance.Stop()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case refund := <-eventChannel:
		s.Raw().NotNil(refund, "refundPegOut failed")
		return RefundPegoutEvent{QuoteHash: hex.EncodeToString(refund.QuoteHash[:]), RawEvent: refund.Raw}
	case err = <-subscription.Err():
		s.Raw().Require().NoError(err)
	case <-done:
		s.Raw().FailNow("Test cancelled while waiting for pegout refund")
	case <-testTolerance.C:
		s.Raw().FailNow("timeout waiting for refund pegout event")
	}
	return RefundPegoutEvent{}
}

func (e *SplitLbcExecutor) GetCallForUserEvent(s TestSuite, timeout time.Duration, userAddress, providerAddress string) CallForUserEvent {
	eventChannel := make(chan *bindings.IPegInCallForUser)
	subscription, err := e.pegin.WatchCallForUser(
		nil,
		eventChannel,
		[]common.Address{common.HexToAddress(providerAddress)},
		[]common.Address{common.HexToAddress(userAddress)},
		nil,
	)
	s.Raw().Require().NoError(err)
	defer subscription.Unsubscribe()

	done := make(chan os.Signal, 1)
	testTolerance := time.NewTimer(timeout)
	defer testTolerance.Stop()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case event := <-eventChannel:
		s.Raw().NotNil(event, "callForUser failed")
		return CallForUserEvent{
			From:      event.From.String(),
			To:        event.Dest.String(),
			QuoteHash: hex.EncodeToString(event.QuoteHash[:]),
			GasLimit:  event.GasLimit.Uint64(),
			Value:     event.Value,
			Data:      event.Data,
			Success:   event.Success,
		}
	case err = <-subscription.Err():
		s.Raw().Require().NoError(err)
	case <-done:
		s.Raw().FailNow("Test cancelled while waiting for call for user")
	case <-testTolerance.C:
		s.Raw().FailNow("timeout waiting for call for user event")
	}
	return CallForUserEvent{}
}

func (e *SplitLbcExecutor) GetPeginRegisteredEvent(s TestSuite, timeout time.Duration, quoteHash string) PegInRegisteredEvent {
	var quoteHashByes [32]byte
	eventChannel := make(chan *bindings.IPegInPegInRegistered)
	parsedHash, err := hex.DecodeString(quoteHash)
	s.Raw().Require().NoError(err)

	copy(quoteHashByes[:], parsedHash)
	subscription, err := e.pegin.WatchPegInRegistered(
		nil,
		eventChannel,
		[][32]byte{quoteHashByes},
		nil,
	)
	s.Raw().Require().NoError(err)
	defer subscription.Unsubscribe()

	done := make(chan os.Signal, 1)
	testTolerance := time.NewTimer(timeout)
	defer testTolerance.Stop()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case event := <-eventChannel:
		s.Raw().NotNil(event, "registerPegin failed")
		return PegInRegisteredEvent{
			QuoteHash: hex.EncodeToString(event.QuoteHash[:]),
			Amount:    entities.NewBigWei(event.TransferredAmount),
			RawEvent:  event.Raw,
		}
	case err = <-subscription.Err():
		s.Raw().Require().NoError(err)
	case <-done:
		s.Raw().FailNow("Test cancelled while waiting for register pegin")
	case <-testTolerance.C:
		s.Raw().FailNow("timeout waiting for register pegin event")
	}
	return PegInRegisteredEvent{}
}

func (e *SplitLbcExecutor) parsePegoutQuote(s *suite.Suite, originalQuote pkg.PegoutQuoteDTO) bindings.QuotesPegOutQuote {
	lpBtcAddress, err := bitcoin.DecodeAddress(originalQuote.LpBTCAddr)
	s.Require().NoError(err)
	btcRefundAddress, err := bitcoin.DecodeAddress(originalQuote.BtcRefundAddr)
	s.Require().NoError(err)
	depositAddress, err := bitcoin.DecodeAddress(originalQuote.DepositAddr)
	s.Require().NoError(err)
	return bindings.QuotesPegOutQuote{
		LbcAddress:            common.HexToAddress(originalQuote.LBCAddr),
		LpRskAddress:          common.HexToAddress(originalQuote.LPRSKAddr),
		BtcRefundAddress:      btcRefundAddress,
		RskRefundAddress:      common.HexToAddress(originalQuote.RSKRefundAddr),
		LpBtcAddress:          lpBtcAddress,
		CallFee:               originalQuote.CallFee,
		PenaltyFee:            originalQuote.PenaltyFee,
		Nonce:                 originalQuote.Nonce,
		DepositAddress:        depositAddress,
		Value:                 originalQuote.Value,
		AgreementTimestamp:    originalQuote.AgreementTimestamp,
		DepositDateLimit:      originalQuote.DepositDateLimit,
		DepositConfirmations:  originalQuote.DepositConfirmations,
		TransferConfirmations: originalQuote.TransferConfirmations,
		TransferTime:          originalQuote.TransferTime,
		ExpireDate:            originalQuote.ExpireDate,
		ExpireBlock:           originalQuote.ExpireBlock,
		GasFee:                originalQuote.GasFee,
	}
}
