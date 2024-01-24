package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"math/rand"
	"strings"
	"time"
)

type GetQuoteUseCase struct {
	rsk                   blockchain.RootstockRpcServer
	feeCollector          blockchain.FeeCollector
	bridge                blockchain.RootstockBridge
	lbc                   blockchain.LiquidityBridgeContract
	pegoutQuoteRepository quote.PegoutQuoteRepository
	lp                    entities.LiquidityProvider
	pegoutLp              entities.PegoutLiquidityProvider
	btcWallet             blockchain.BitcoinWallet
	feeCollectorAddress   string
}

func NewGetQuoteUseCase(
	rsk blockchain.RootstockRpcServer,
	feeCollector blockchain.FeeCollector,
	bridge blockchain.RootstockBridge,
	lbc blockchain.LiquidityBridgeContract,
	pegoutQuoteRepository quote.PegoutQuoteRepository,
	lp entities.LiquidityProvider,
	pegoutLp entities.PegoutLiquidityProvider,
	btcWallet blockchain.BitcoinWallet,
	feeCollectorAddress string,
) *GetQuoteUseCase {
	return &GetQuoteUseCase{
		rsk:                   rsk,
		feeCollector:          feeCollector,
		bridge:                bridge,
		lbc:                   lbc,
		pegoutQuoteRepository: pegoutQuoteRepository,
		lp:                    lp,
		pegoutLp:              pegoutLp,
		btcWallet:             btcWallet,
		feeCollectorAddress:   feeCollectorAddress,
	}
}

type QuoteRequest struct {
	to                   string
	valueToTransfer      *entities.Wei
	rskRefundAddress     string
	bitcoinRefundAddress string
}

func NewQuoteRequest(
	to string,
	valueToTransfer *entities.Wei,
	rskRefundAddress string,
	bitcoinRefundAddress string,
) QuoteRequest {
	return QuoteRequest{
		to:                   to,
		valueToTransfer:      valueToTransfer,
		rskRefundAddress:     rskRefundAddress,
		bitcoinRefundAddress: bitcoinRefundAddress,
	}
}

type GetPegoutQuoteResult struct {
	PegoutQuote quote.PegoutQuote
	Hash        string
}

func (useCase *GetQuoteUseCase) Run(ctx context.Context, request QuoteRequest) (GetPegoutQuoteResult, error) {
	var daoFeePercentage, blockNumber uint64
	var daoTxAmounts usecases.DaoAmounts
	var hash string

	feeInWei := new(entities.Wei)
	gasPrice := new(entities.Wei)
	gasFeeDao := new(entities.Wei)
	minLockTxValueInSatoshi := new(entities.Wei)
	errorArgs := usecases.NewErrorArgs()
	var err error

	if errorArgs, err = useCase.validateRequest(request); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPegoutQuoteId, err, errorArgs)
	}

	if feeInWei, err = useCase.btcWallet.EstimateTxFees(request.to, request.valueToTransfer); err != nil && strings.Contains(err.Error(), "Insufficient Funds") {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, usecases.NoLiquidityError)
	} else if err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if daoFeePercentage, err = useCase.feeCollector.DaoFeePercentage(); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if daoTxAmounts, err = usecases.CalculateDaoAmounts(ctx, useCase.rsk, request.valueToTransfer, daoFeePercentage, useCase.feeCollectorAddress); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if blockNumber, err = useCase.rsk.GetHeight(ctx); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if gasPrice, err = useCase.rsk.GasPrice(ctx); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	gasFeeDao.Mul(daoTxAmounts.DaoGasAmount, gasPrice)
	totalGasFee := new(entities.Wei).Add(feeInWei, gasFeeDao)
	now := uint32(time.Now().Unix())
	confirmationsForUserTx := useCase.lp.GetRootstockConfirmationsForValue(request.valueToTransfer)
	confirmationsForLpTx := useCase.lp.GetBitcoinConfirmationsForValue(request.valueToTransfer)
	pegoutQuote := quote.PegoutQuote{
		LbcAddress:            useCase.lbc.GetAddress(),
		LpRskAddress:          useCase.lp.RskAddress(),
		BtcRefundAddress:      request.bitcoinRefundAddress,
		RskRefundAddress:      request.rskRefundAddress,
		LpBtcAddress:          useCase.lp.BtcAddress(),
		CallFee:               useCase.pegoutLp.CallFeePegout(),
		PenaltyFee:            useCase.pegoutLp.PenaltyFeePegout().Uint64(),
		Nonce:                 int64(rand.Int()),
		DepositAddress:        request.to,
		Value:                 request.valueToTransfer,
		AgreementTimestamp:    now,
		DepositDateLimit:      now + useCase.pegoutLp.TimeForDepositPegout(),
		DepositConfirmations:  confirmationsForUserTx,
		TransferConfirmations: confirmationsForLpTx,
		TransferTime:          useCase.pegoutLp.TimeForDepositPegout(),
		ExpireDate:            now + useCase.pegoutLp.TimeForDepositPegout(),
		ExpireBlock:           uint32(blockNumber + useCase.pegoutLp.ExpireBlocksPegout()),
		GasFee:                totalGasFee,
		ProductFeeAmount:      daoTxAmounts.DaoFeeAmount.Uint64(),
	}

	if err = entities.ValidateStruct(pegoutQuote); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if minLockTxValueInSatoshi, err = useCase.bridge.GetMinimumLockTxValue(); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	minimumInWei := entities.SatoshiToWei(minLockTxValueInSatoshi.Uint64())
	if pegoutQuote.Total().Cmp(minimumInWei) <= 0 {
		errorArgs["minimum"] = minimumInWei.String()
		errorArgs["value"] = pegoutQuote.Total().String()
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPegoutQuoteId, usecases.TxBelowMinimumError, errorArgs)
	}

	if hash, err = useCase.lbc.HashPegoutQuote(pegoutQuote); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if err = useCase.pegoutQuoteRepository.InsertQuote(ctx, hash, pegoutQuote); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return GetPegoutQuoteResult{
		PegoutQuote: pegoutQuote,
		Hash:        hash,
	}, nil
}

func (useCase *GetQuoteUseCase) validateRequest(request QuoteRequest) (usecases.ErrorArgs, error) {
	var err error
	errorArgs := usecases.NewErrorArgs()
	if !blockchain.IsLegacyBtcAddress(request.to) {
		errorArgs["btcAddress"] = request.to
		return errorArgs, usecases.BtcAddressNotSupportedError
	} else if !blockchain.IsLegacyBtcAddress(request.bitcoinRefundAddress) {
		errorArgs["btcAddress"] = request.bitcoinRefundAddress
		return errorArgs, usecases.BtcAddressNotSupportedError
	} else if !blockchain.IsRskAddress(request.rskRefundAddress) {
		errorArgs["rskAddress"] = request.rskRefundAddress
		return errorArgs, usecases.RskAddressNotSupportedError
	} else if err = useCase.pegoutLp.ValidateAmountForPegout(request.valueToTransfer); err != nil {
		return errorArgs, err
	} else {
		return nil, nil
	}
}
