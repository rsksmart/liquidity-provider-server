package pegin

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"math/rand"
	"time"
)

type GetQuoteUseCase struct {
	rsk                  blockchain.RootstockRpcServer
	feeCollector         blockchain.FeeCollector
	bridge               blockchain.RootstockBridge
	lbc                  blockchain.LiquidityBridgeContract
	peginQuoteRepository quote.PeginQuoteRepository
	lp                   entities.LiquidityProvider
	peginLp              entities.PeginLiquidityProvider
	feeCollectorAddress  string
}

func NewGetQuoteUseCase(
	rsk blockchain.RootstockRpcServer,
	feeCollector blockchain.FeeCollector,
	bridge blockchain.RootstockBridge,
	lbc blockchain.LiquidityBridgeContract,
	peginQuoteRepository quote.PeginQuoteRepository,
	lp entities.LiquidityProvider,
	peginLp entities.PeginLiquidityProvider,
	feeCollectorAddress string,
) *GetQuoteUseCase {
	return &GetQuoteUseCase{
		rsk:                  rsk,
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
	var fedAddress, hash string
	var daoFeePercentage uint64
	var errorArgs usecases.ErrorArgs
	var err error
	var gasPrice, estimatedCallGas *entities.Wei

	minLockTxValueInSatoshi := new(entities.Wei)

	if errorArgs, err = useCase.validateRequest(request); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPeginQuoteId, err, errorArgs)
	}

	estimatedCallGas, err = useCase.rsk.EstimateGas(ctx, request.callEoaOrContractAddress, request.valueToTransfer, request.callContractArguments)
	if err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if gasPrice, err = useCase.rsk.GasPrice(ctx); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if daoFeePercentage, err = useCase.feeCollector.DaoFeePercentage(); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	if daoTxAmounts, err = usecases.CalculateDaoAmounts(ctx, useCase.rsk, request.valueToTransfer, daoFeePercentage, useCase.feeCollectorAddress); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if fedAddress, err = useCase.bridge.GetFedAddress(); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	totalGas := new(entities.Wei).Add(estimatedCallGas, daoTxAmounts.DaoGasAmount)
	gasFee := new(entities.Wei).Mul(totalGas, gasPrice)
	peginQuote := quote.PeginQuote{
		FedBtcAddress:      fedAddress,
		LbcAddress:         useCase.lbc.GetAddress(),
		LpRskAddress:       useCase.lp.RskAddress(),
		BtcRefundAddress:   request.bitcoinRefundAddress,
		RskRefundAddress:   request.rskRefundAddress,
		LpBtcAddress:       useCase.lp.BtcAddress(),
		CallFee:            useCase.peginLp.CallFeePegin(),
		PenaltyFee:         useCase.peginLp.PenaltyFeePegin(),
		ContractAddress:    request.callEoaOrContractAddress,
		Data:               hex.EncodeToString(request.callContractArguments),
		GasLimit:           uint32(totalGas.Uint64()),
		Nonce:              int64(rand.Int()),
		Value:              request.valueToTransfer,
		AgreementTimestamp: uint32(time.Now().Unix()),
		TimeForDeposit:     useCase.peginLp.TimeForDepositPegin(),
		LpCallTime:         useCase.peginLp.CallTime(),
		Confirmations:      useCase.lp.GetBitcoinConfirmationsForValue(request.valueToTransfer),
		CallOnRegister:     false,
		GasFee:             gasFee,
		ProductFeeAmount:   daoTxAmounts.DaoFeeAmount.Uint64(),
	}

	if err = entities.ValidateStruct(peginQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if minLockTxValueInSatoshi, err = useCase.bridge.GetMinimumLockTxValue(); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	minimumInWei := entities.SatoshiToWei(minLockTxValueInSatoshi.Uint64())
	if peginQuote.Total().Cmp(minimumInWei) <= 0 {
		errorArgs["minimum"] = minimumInWei.String()
		errorArgs["value"] = peginQuote.Total().String()
		return GetPeginQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPeginQuoteId, usecases.TxBelowMinimumError, errorArgs)
	}

	if hash, err = useCase.lbc.HashPeginQuote(peginQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if err = useCase.peginQuoteRepository.InsertQuote(ctx, hash, peginQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	return GetPeginQuoteResult{
		PeginQuote: peginQuote,
		Hash:       hash,
	}, nil
}

func (useCase *GetQuoteUseCase) validateRequest(request QuoteRequest) (usecases.ErrorArgs, error) {
	var err error
	args := usecases.NewErrorArgs()
	if !blockchain.IsLegacyBtcAddress(request.bitcoinRefundAddress) {
		args["btcAddress"] = request.bitcoinRefundAddress
		return args, usecases.BtcAddressNotSupportedError
	} else if !blockchain.IsRskAddress(request.rskRefundAddress) {
		args["rskAddress"] = request.rskRefundAddress
		return args, usecases.RskAddressNotSupportedError
	} else if !blockchain.IsRskAddress(request.callEoaOrContractAddress) {
		args["rskAddress"] = request.callEoaOrContractAddress
		return args, usecases.RskAddressNotSupportedError
	} else if err = useCase.peginLp.ValidateAmountForPegin(request.valueToTransfer); err != nil {
		return args, err
	} else {
		return nil, nil
	}
}
