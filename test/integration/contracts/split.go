package contracts

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	collateralBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	discoveryBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	pegoutBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/integration"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Contract[binding any] struct {
	contract *bind.BoundContract
	binding  *binding
}

type SplitLbcExecutor struct {
	pegout               *Contract[pegoutBindings.PegoutContract]
	discovery            *Contract[discoveryBindings.FlyoverDiscovery]
	collateralManagement *Contract[collateralBindings.CollateralManagementContract]
}

type SplitAddresses struct {
	Discovery            string
	Pegout               string
	Pegin                string
	CollateralManagement string
}

func NewSplityLbcExecutor(addresses SplitAddresses, backend bind.ContractBackend) (*SplitLbcExecutor, error) {
	pegoutBinding := pegoutBindings.NewPegoutContract()
	discoveryBinding := discoveryBindings.NewFlyoverDiscovery()
	collateralBinding := collateralBindings.NewCollateralManagementContract()
	pegout := &Contract[pegoutBindings.PegoutContract]{
		contract: pegoutBinding.Instance(backend, common.HexToAddress(addresses.Pegout)),
		binding:  pegoutBinding,
	}
	discovery := &Contract[discoveryBindings.FlyoverDiscovery]{
		contract: discoveryBinding.Instance(backend, common.HexToAddress(addresses.Discovery)),
		binding:  discoveryBinding,
	}
	collateralManagement := &Contract[collateralBindings.CollateralManagementContract]{
		contract: collateralBinding.Instance(backend, common.HexToAddress(addresses.CollateralManagement)),
		binding:  collateralBinding,
	}
	return &SplitLbcExecutor{pegout: pegout, discovery: discovery, collateralManagement: collateralManagement}, nil
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

	depositTx, err := bind.Transact(e.pegout.contract, opts, e.pegout.binding.PackDepositPegOut(parsedQuote, signature))
	s.Raw().Require().NoError(err)
	receipt, err := bind.WaitMined(ctx, s.RskClient(), depositTx.Hash())
	s.Raw().Require().NoError(err)
	log.Info("[Integration test] Hash of deposit tx ", depositTx.Hash().String())
	return receipt, depositTx
}

func (e *SplitLbcExecutor) GetRefundPegoutEvent(s TestSuite, timeout time.Duration, quoteHash string) RefundPegoutEvent {
	var quoteHashByes [32]byte
	eventChannel := make(chan *pegoutBindings.PegoutContractPegOutRefunded)
	parsedHash, err := hex.DecodeString(quoteHash)
	s.Raw().Require().NoError(err)
	copy(quoteHashByes[:], parsedHash)

	subscription, err := bind.WatchEvents(e.pegout.contract, &bind.WatchOpts{}, e.pegout.binding.UnpackPegOutRefundedEvent, eventChannel, []any{quoteHashByes})
	s.Raw().Require().NoError(err)
	defer subscription.Unsubscribe()

	done := make(chan os.Signal, 1)
	testTolerance := time.NewTimer(timeout)
	defer testTolerance.Stop()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case refund := <-eventChannel:
		s.Raw().NotNil(refund, "refundPegOut failed")
		return RefundPegoutEvent{QuoteHash: hex.EncodeToString(refund.QuoteHash[:]), RawEvent: *refund.Raw}
	case err = <-subscription.Err():
		s.Raw().Require().NoError(err)
	case <-done:
		s.Raw().FailNow("Test cancelled while waiting for pegout refund")
	case <-testTolerance.C:
		s.Raw().FailNow("timeout waiting for refund pegout event")
	}
	return RefundPegoutEvent{}
}

func (e *SplitLbcExecutor) parsePegoutQuote(s *suite.Suite, originalQuote pkg.PegoutQuoteDTO) pegoutBindings.QuotesPegOutQuote {
	lpBtcAddress, err := bitcoin.DecodeAddress(originalQuote.LpBTCAddr)
	s.Require().NoError(err)
	btcRefundAddress, err := bitcoin.DecodeAddress(originalQuote.BtcRefundAddr)
	s.Require().NoError(err)
	depositAddress, err := bitcoin.DecodeAddress(originalQuote.DepositAddr)
	s.Require().NoError(err)
	return pegoutBindings.QuotesPegOutQuote{
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
