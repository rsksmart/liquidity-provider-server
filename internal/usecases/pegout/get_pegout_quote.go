package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
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
	var pegoutQuote quote.PegoutQuote
	var daoTxAmounts usecases.DaoAmounts
	var hash string
	var errorArgs usecases.ErrorArgs
	var gasPrice, feeInWei *entities.Wei
	var err error

	gasFeeDao := new(entities.Wei)
	if errorArgs, err = useCase.validateRequest(request); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPegoutQuoteId, err, errorArgs)
	}

	if feeInWei, err = useCase.btcWallet.EstimateTxFees(request.to, request.valueToTransfer); err != nil && strings.Contains(err.Error(), "Insufficient Funds") {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, usecases.NoLiquidityError)
	} else if err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if daoTxAmounts, err = useCase.buildDaoAmounts(ctx, request); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	if gasPrice, err = useCase.rsk.GasPrice(ctx); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	gasFeeDao.Mul(daoTxAmounts.DaoGasAmount, gasPrice)
	fees := quote.Fees{
		CallFee:          useCase.pegoutLp.CallFeePegout(),
		GasFee:           new(entities.Wei).Add(feeInWei, gasFeeDao),
		PenaltyFee:       useCase.pegoutLp.PenaltyFeePegout(),
		ProductFeeAmount: daoTxAmounts.DaoFeeAmount.Uint64(),
	}
	if pegoutQuote, err = useCase.buildPegoutQuote(ctx, request, fees); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	if err = usecases.ValidateMinLockValue(usecases.GetPegoutQuoteId, useCase.bridge, pegoutQuote.Total()); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	if hash, err = useCase.persistQuote(ctx, pegoutQuote); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	return GetPegoutQuoteResult{PegoutQuote: pegoutQuote, Hash: hash}, nil
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

func (useCase *GetQuoteUseCase) buildPegoutQuote(ctx context.Context, request QuoteRequest, fees quote.Fees) (quote.PegoutQuote, error) {
	var err error
	var nonce int64
	var blockNumber uint64

	if nonce, err = usecases.GetRandomInt(); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if blockNumber, err = useCase.rsk.GetHeight(ctx); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	now := uint32(time.Now().Unix())
	confirmationsForUserTx := useCase.lp.GetRootstockConfirmationsForValue(request.valueToTransfer)
	confirmationsForLpTx := useCase.lp.GetBitcoinConfirmationsForValue(request.valueToTransfer)
	pegoutQuote := quote.PegoutQuote{
		LbcAddress:            useCase.lbc.GetAddress(),
		LpRskAddress:          useCase.lp.RskAddress(),
		BtcRefundAddress:      request.bitcoinRefundAddress,
		RskRefundAddress:      request.rskRefundAddress,
		LpBtcAddress:          useCase.lp.BtcAddress(),
		CallFee:               fees.CallFee,
		PenaltyFee:            fees.PenaltyFee.Uint64(),
		Nonce:                 nonce,
		DepositAddress:        request.to,
		Value:                 request.valueToTransfer,
		AgreementTimestamp:    now,
		DepositDateLimit:      now + useCase.pegoutLp.TimeForDepositPegout(),
		DepositConfirmations:  confirmationsForUserTx,
		TransferConfirmations: confirmationsForLpTx,
		TransferTime:          useCase.pegoutLp.TimeForDepositPegout(),
		ExpireDate:            now + useCase.pegoutLp.TimeForDepositPegout(),
		ExpireBlock:           uint32(blockNumber + useCase.pegoutLp.ExpireBlocksPegout()),
		GasFee:                fees.GasFee,
		ProductFeeAmount:      fees.ProductFeeAmount,
	}

	if err = entities.ValidateStruct(pegoutQuote); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return pegoutQuote, nil
}

func (useCase *GetQuoteUseCase) buildDaoAmounts(ctx context.Context, request QuoteRequest) (usecases.DaoAmounts, error) {
	var daoTxAmounts usecases.DaoAmounts
	var daoFeePercentage uint64
	var err error
	if daoFeePercentage, err = useCase.feeCollector.DaoFeePercentage(); err != nil {
		return usecases.DaoAmounts{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	if daoTxAmounts, err = usecases.CalculateDaoAmounts(ctx, useCase.rsk, request.valueToTransfer, daoFeePercentage, useCase.feeCollectorAddress); err != nil {
		return usecases.DaoAmounts{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return daoTxAmounts, nil
}

func (useCase *GetQuoteUseCase) persistQuote(ctx context.Context, pegoutQuote quote.PegoutQuote) (string, error) {
	var hash string
	var err error
	if hash, err = useCase.lbc.HashPegoutQuote(pegoutQuote); err != nil {
		return "", usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if err = useCase.pegoutQuoteRepository.InsertQuote(ctx, hash, pegoutQuote); err != nil {
		return "", usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return hash, nil
}
