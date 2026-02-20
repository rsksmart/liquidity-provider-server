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
}

func NewGetQuoteUseCase(
	rpc blockchain.Rpc,
	contracts blockchain.RskContracts,
	pegoutQuoteRepository quote.PegoutQuoteRepository,
	lp liquidity_provider.LiquidityProvider,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	btcWallet blockchain.BitcoinWallet,
) *GetQuoteUseCase {
	return &GetQuoteUseCase{
		rpc:                   rpc,
		contracts:             contracts,
		pegoutQuoteRepository: pegoutQuoteRepository,
		lp:                    lp,
		pegoutLp:              pegoutLp,
		btcWallet:             btcWallet,
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
	var hash string
	var errorArgs usecases.ErrorArgs
	var btcFeeEstimation blockchain.BtcFeeEstimation
	var creationData quote.PegoutCreationData
	var err error

	if err = usecases.CheckPauseState(useCase.contracts.PegOut); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	configuration := useCase.pegoutLp.PegoutConfiguration(ctx)
	if errorArgs, err = useCase.validateRequest(configuration, request); err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPegoutQuoteId, err, errorArgs)
	}

	if btcFeeEstimation, err = useCase.btcWallet.EstimateTxFees(request.to, request.valueToTransfer); err != nil &&
		strings.Contains(strings.ToLower(err.Error()), "insufficient funds") {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, usecases.NoLiquidityError)
	} else if err != nil {
		return GetPegoutQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	if creationData, err = useCase.buildCreationData(ctx, btcFeeEstimation, configuration); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	fees := quote.Fees{
		CallFee:    quote.CalculateCallFee(request.valueToTransfer, configuration),
		GasFee:     btcFeeEstimation.Value,
		PenaltyFee: configuration.PenaltyFee,
	}
	if pegoutQuote, err = useCase.buildPegoutQuote(ctx, configuration, request, fees); err != nil {
		return GetPegoutQuoteResult{}, err
	}

	if hash, err = useCase.persistQuote(ctx, pegoutQuote, creationData); err != nil {
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

	chainId, err := useCase.rpc.Rsk.ChainId(ctx)
	if err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	now := uint32(time.Now().Unix())
	generalConfiguration := useCase.lp.GeneralConfiguration(ctx)
	confirmationsForUserTx := generalConfiguration.RskConfirmations.ForValue(request.valueToTransfer)
	confirmationsForLpTx := generalConfiguration.BtcConfirmations.ForValue(request.valueToTransfer)
	pegoutQuote := quote.PegoutQuote{
		LbcAddress:            useCase.contracts.PegOut.GetAddress(),
		LpRskAddress:          useCase.lp.RskAddress(),
		BtcRefundAddress:      request.to,
		RskRefundAddress:      request.rskRefundAddress,
		LpBtcAddress:          useCase.lp.BtcAddress(),
		CallFee:               fees.CallFee,
		PenaltyFee:            fees.PenaltyFee,
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
		ChainId:               chainId,
	}

	if err = entities.ValidateStruct(pegoutQuote); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return pegoutQuote, nil
}

func (useCase *GetQuoteUseCase) persistQuote(ctx context.Context, pegoutQuote quote.PegoutQuote, creationData quote.PegoutCreationData) (string, error) {
	var hash string
	var err error
	if hash, err = useCase.contracts.PegOut.HashPegoutQuote(pegoutQuote); err != nil {
		return "", usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}

	createdQuote := quote.CreatedPegoutQuote{Quote: pegoutQuote, CreationData: creationData, Hash: hash}
	if err = useCase.pegoutQuoteRepository.InsertQuote(ctx, createdQuote); err != nil {
		return "", usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	return hash, nil
}

func (useCase *GetQuoteUseCase) buildCreationData(
	ctx context.Context,
	btcFeeEstimation blockchain.BtcFeeEstimation,
	configuration liquidity_provider.PegoutConfiguration,
) (quote.PegoutCreationData, error) {
	var gasPrice *entities.Wei
	var err error

	if gasPrice, err = useCase.rpc.Rsk.GasPrice(ctx); err != nil {
		return quote.PegoutCreationData{}, usecases.WrapUseCaseError(usecases.GetPegoutQuoteId, err)
	}
	creationData := quote.PegoutCreationData{
		FeeRate:       btcFeeEstimation.FeeRate,
		GasPrice:      gasPrice,
		FeePercentage: configuration.FeePercentage,
		FixedFee:      configuration.FixedFee,
	}
	return creationData, nil
}
