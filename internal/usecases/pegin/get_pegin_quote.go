package pegin

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"time"
)

type GetQuoteUseCase struct {
	rsk                  blockchain.RootstockRpcServer
	btc                  blockchain.BitcoinNetwork
	feeCollector         blockchain.FeeCollector
	bridge               blockchain.RootstockBridge
	lbc                  blockchain.LiquidityBridgeContract
	peginQuoteRepository quote.PeginQuoteRepository
	lp                   liquidity_provider.LiquidityProvider
	peginLp              liquidity_provider.PeginLiquidityProvider
	feeCollectorAddress  string
}

func NewGetQuoteUseCase(
	rsk blockchain.RootstockRpcServer,
	btc blockchain.BitcoinNetwork,
	feeCollector blockchain.FeeCollector,
	bridge blockchain.RootstockBridge,
	lbc blockchain.LiquidityBridgeContract,
	peginQuoteRepository quote.PeginQuoteRepository,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	feeCollectorAddress string,
) *GetQuoteUseCase {
	return &GetQuoteUseCase{
		rsk:                  rsk,
		btc:                  btc,
		feeCollector:         feeCollector,
		bridge:               bridge,
		lbc:                  lbc,
		peginQuoteRepository: peginQuoteRepository,
		lp:                   lp,
		peginLp:              peginLp,
		feeCollectorAddress:  feeCollectorAddress,
	}
}

type QuoteRequest struct {
	callEoaOrContractAddress string
	callContractArguments    []byte
	valueToTransfer          *entities.Wei
	rskRefundAddress         string
	bitcoinRefundAddress     string
}

func NewQuoteRequest(
	callEoaOrContractAddress string,
	callContractArguments []byte,
	valueToTransfer *entities.Wei,
	rskRefundAddress string,
	bitcoinRefundAddress string,
) QuoteRequest {
	return QuoteRequest{
		callEoaOrContractAddress: callEoaOrContractAddress,
		callContractArguments:    callContractArguments,
		valueToTransfer:          valueToTransfer,
		rskRefundAddress:         rskRefundAddress,
		bitcoinRefundAddress:     bitcoinRefundAddress,
	}
}

type GetPeginQuoteResult struct {
	PeginQuote quote.PeginQuote
	Hash       string
}

func (useCase *GetQuoteUseCase) Run(ctx context.Context, request QuoteRequest) (GetPeginQuoteResult, error) {
	var daoTxAmounts usecases.DaoAmounts
	var peginQuote quote.PeginQuote
	var fedAddress, hash string
	var errorArgs usecases.ErrorArgs
	var err error
	var gasPrice, estimatedCallGas *entities.Wei

	peginConfiguration := useCase.peginLp.PeginConfiguration(ctx)
	if errorArgs, err = useCase.validateRequest(peginConfiguration, request); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPeginQuoteId, err, errorArgs)
	}

	estimatedCallGas, err = useCase.rsk.EstimateGas(ctx, request.callEoaOrContractAddress, request.valueToTransfer, request.callContractArguments)
	if err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if gasPrice, err = useCase.rsk.GasPrice(ctx); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if daoTxAmounts, err = useCase.buildDaoAmounts(ctx, request); err != nil {
		return GetPeginQuoteResult{}, err
	}

	if fedAddress, err = useCase.bridge.GetFedAddress(); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	generalConfiguration := useCase.lp.GeneralConfiguration(ctx)
	totalGas := new(entities.Wei).Add(estimatedCallGas, daoTxAmounts.DaoGasAmount)
	fees := quote.Fees{
		CallFee:          peginConfiguration.CallFee,
		GasFee:           new(entities.Wei).Mul(totalGas, gasPrice),
		PenaltyFee:       peginConfiguration.PenaltyFee,
		ProductFeeAmount: daoTxAmounts.DaoFeeAmount.Uint64(),
	}
	if peginQuote, err = useCase.buildPeginQuote(generalConfiguration, peginConfiguration, request, fedAddress, totalGas, fees); err != nil {
		return GetPeginQuoteResult{}, err
	}

	if err = usecases.ValidateMinLockValue(usecases.GetPeginQuoteId, useCase.bridge, peginQuote.Total()); err != nil {
		return GetPeginQuoteResult{}, err
	}

	if hash, err = useCase.lbc.HashPeginQuote(peginQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	if err = useCase.peginQuoteRepository.InsertQuote(ctx, hash, peginQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	return GetPeginQuoteResult{PeginQuote: peginQuote, Hash: hash}, nil
}

func (useCase *GetQuoteUseCase) validateRequest(configuration liquidity_provider.PeginConfiguration, request QuoteRequest) (usecases.ErrorArgs, error) {
	var err error
	args := usecases.NewErrorArgs()
	if err = useCase.btc.ValidateAddress(request.bitcoinRefundAddress); err != nil {
		args["btcAddress"] = request.bitcoinRefundAddress
		return args, err
	}
	if !blockchain.IsRskAddress(request.rskRefundAddress) {
		args["rskAddress"] = request.rskRefundAddress
		return args, usecases.RskAddressNotSupportedError
	}
	if !blockchain.IsRskAddress(request.callEoaOrContractAddress) {
		args["rskAddress"] = request.callEoaOrContractAddress
		return args, usecases.RskAddressNotSupportedError
	}
	if err = configuration.ValidateAmount(request.valueToTransfer); err != nil {
		return args, err
	}
	return nil, nil
}

func (useCase *GetQuoteUseCase) buildPeginQuote(
	generalConfig liquidity_provider.GeneralConfiguration,
	peginConfig liquidity_provider.PeginConfiguration,
	request QuoteRequest,
	fedAddress string,
	totalGas *entities.Wei,
	fees quote.Fees,
) (quote.PeginQuote, error) {
	var err error
	var nonce int64

	if nonce, err = usecases.GetRandomInt(); err != nil {
		return quote.PeginQuote{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	peginQuote := quote.PeginQuote{
		FedBtcAddress:      fedAddress,
		LbcAddress:         useCase.lbc.GetAddress(),
		LpRskAddress:       useCase.lp.RskAddress(),
		BtcRefundAddress:   request.bitcoinRefundAddress,
		RskRefundAddress:   request.rskRefundAddress,
		LpBtcAddress:       useCase.lp.BtcAddress(),
		CallFee:            fees.CallFee,
		PenaltyFee:         fees.PenaltyFee,
		ContractAddress:    request.callEoaOrContractAddress,
		Data:               hex.EncodeToString(request.callContractArguments),
		GasLimit:           uint32(totalGas.Uint64()),
		Nonce:              nonce,
		Value:              request.valueToTransfer,
		AgreementTimestamp: uint32(time.Now().Unix()),
		TimeForDeposit:     peginConfig.TimeForDeposit,
		LpCallTime:         peginConfig.CallTime,
		Confirmations:      generalConfig.BtcConfirmations.ForValue(request.valueToTransfer),
		CallOnRegister:     false,
		GasFee:             fees.GasFee,
		ProductFeeAmount:   fees.ProductFeeAmount,
	}

	if err = entities.ValidateStruct(peginQuote); err != nil {
		return quote.PeginQuote{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	return peginQuote, nil
}

func (useCase *GetQuoteUseCase) buildDaoAmounts(ctx context.Context, request QuoteRequest) (usecases.DaoAmounts, error) {
	var daoTxAmounts usecases.DaoAmounts
	var daoFeePercentage uint64
	var err error

	if daoFeePercentage, err = useCase.feeCollector.DaoFeePercentage(); err != nil {
		return usecases.DaoAmounts{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	if daoTxAmounts, err = usecases.CalculateDaoAmounts(ctx, useCase.rsk, request.valueToTransfer, daoFeePercentage, useCase.feeCollectorAddress); err != nil {
		return usecases.DaoAmounts{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	return daoTxAmounts, nil
}
