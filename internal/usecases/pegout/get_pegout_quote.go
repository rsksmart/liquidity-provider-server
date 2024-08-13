package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"strings"
	"time"
)

type GetQuoteUseCase struct {
	rpc                   blockchain.Rpc
	contracts             blockchain.RskContracts
	pegoutQuoteRepository quote.PegoutQuoteRepository
	lp                    liquidity_provider.LiquidityProvider
	pegoutLp              liquidity_provider.PegoutLiquidityProvider
	btcWallet             blockchain.BitcoinWallet
	feeCollectorAddress   string
}

func NewGetQuoteUseCase(
	rpc blockchain.Rpc,
	contracts blockchain.RskContracts,
	pegoutQuoteRepository quote.PegoutQuoteRepository,
	lp liquidity_provider.LiquidityProvider,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	btcWallet blockchain.BitcoinWallet,
	feeCollectorAddress string,
) *GetQuoteUseCase {
	return &GetQuoteUseCase{
		rpc:                   rpc,
		contracts:             contracts,
		pegoutQuoteRepository: pegoutQuoteRepository,
		lp:                    lp,
		pegoutLp:              pegoutLp,
		btcWallet:             btcWallet,
		feeCollectorAddress:   feeCollectorAddress,
	}
}

type QuoteRequest struct {
	to               string
	valueToTransfer  *entities.Wei
	rskRefundAddress string
}

func NewQuoteRequest(
	to string,
	valueToTransfer *entities.Wei,
	rskRefundAddress string,
) QuoteRequest {
	return QuoteRequest{
		to:               to,
		valueToTransfer:  valueToTransfer,
		rskRefundAddress: rskRefundAddress,
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
	configuration := useCase.pegoutLp.PegoutConfiguration(ctx)
	if errorArgs, err = useCase.validateRequest(configuration, request); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPegoutQuoteId, err, errorArgs)
	}

	if feeInWei, err = useCase.btcWallet.EstimateTxFees(request.to, request.valueToTransfer); err != nil &&
		strings.Contains(strings.ToLower(err.Error()), "insufficient funds") {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, usecases.NoLiquidityError)
	} else if err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if daoTxAmounts, err = useCase.buildDaoAmounts(ctx, request); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	if gasPrice, err = useCase.rpc.Rsk.GasPrice(ctx); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	gasFeeDao.Mul(daoTxAmounts.DaoGasAmount, gasPrice)
	fees := quote.Fees{
		CallFee:          configuration.CallFee,
		GasFee:           new(entities.Wei).Add(feeInWei, gasFeeDao),
		PenaltyFee:       configuration.PenaltyFee,
		ProductFeeAmount: daoTxAmounts.DaoFeeAmount.Uint64(),
	}
	if pegoutQuote, err = useCase.buildPegoutQuote(ctx, configuration, request, fees); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	if hash, err = useCase.persistQuote(ctx, pegoutQuote); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	return GetPegoutQuoteResult{PegoutQuote: pegoutQuote, Hash: hash}, nil
}

func (useCase *GetQuoteUseCase) validateRequest(configuration liquidity_provider.PegoutConfiguration, request QuoteRequest) (usecases.ErrorArgs, error) {
	var err error
	errorArgs := usecases.NewErrorArgs()
	if err = useCase.rpc.Btc.ValidateAddress(request.to); err != nil {
		errorArgs["btcAddress"] = request.to
		return errorArgs, err
	}
	if !blockchain.IsRskAddress(request.rskRefundAddress) {
		errorArgs["rskAddress"] = request.rskRefundAddress
		return errorArgs, usecases.RskAddressNotSupportedError
	}
	if err = configuration.ValidateAmount(request.valueToTransfer); err != nil {
		return errorArgs, err
	}
	return nil, nil
}

func (useCase *GetQuoteUseCase) buildPegoutQuote(
	ctx context.Context,
	configuration liquidity_provider.PegoutConfiguration,
	request QuoteRequest,
	fees quote.Fees,
) (quote.PegoutQuote, error) {
	var err error
	var nonce int64
	var blockNumber uint64

	if nonce, err = utils.GetRandomInt(); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if blockNumber, err = useCase.rpc.Rsk.GetHeight(ctx); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	now := uint32(time.Now().Unix())
	generalConfiguration := useCase.lp.GeneralConfiguration(ctx)
	confirmationsForUserTx := generalConfiguration.RskConfirmations.ForValue(request.valueToTransfer)
	confirmationsForLpTx := generalConfiguration.BtcConfirmations.ForValue(request.valueToTransfer)
	pegoutQuote := quote.PegoutQuote{
		LbcAddress:            useCase.contracts.Lbc.GetAddress(),
		LpRskAddress:          useCase.lp.RskAddress(),
		BtcRefundAddress:      request.to,
		RskRefundAddress:      request.rskRefundAddress,
		LpBtcAddress:          useCase.lp.BtcAddress(),
		CallFee:               fees.CallFee,
		PenaltyFee:            fees.PenaltyFee.Uint64(),
		Nonce:                 nonce,
		DepositAddress:        request.to,
		Value:                 request.valueToTransfer,
		AgreementTimestamp:    now,
		DepositDateLimit:      now + configuration.TimeForDeposit,
		DepositConfirmations:  confirmationsForUserTx,
		TransferConfirmations: confirmationsForLpTx,
		TransferTime:          configuration.TimeForDeposit,
		ExpireDate:            now + configuration.ExpireTime,
		ExpireBlock:           uint32(blockNumber + configuration.ExpireBlocks),
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
	if daoFeePercentage, err = useCase.contracts.FeeCollector.DaoFeePercentage(); err != nil {
		return usecases.DaoAmounts{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	if daoTxAmounts, err = usecases.CalculateDaoAmounts(ctx, useCase.rpc.Rsk, request.valueToTransfer, daoFeePercentage, useCase.feeCollectorAddress); err != nil {
		return usecases.DaoAmounts{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return daoTxAmounts, nil
}

func (useCase *GetQuoteUseCase) persistQuote(ctx context.Context, pegoutQuote quote.PegoutQuote) (string, error) {
	var hash string
	var err error
	if hash, err = useCase.contracts.Lbc.HashPegoutQuote(pegoutQuote); err != nil {
		return "", usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if err = useCase.pegoutQuoteRepository.InsertQuote(ctx, hash, pegoutQuote); err != nil {
		return "", usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return hash, nil
}
