package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"time"
)

type GetQuoteUseCase struct {
	rpc                  blockchain.Rpc
	contracts            blockchain.RskContracts
	peginQuoteRepository quote.PeginQuoteRepository
	lp                   liquidity_provider.LiquidityProvider
	peginLp              liquidity_provider.PeginLiquidityProvider
}

func NewGetQuoteUseCase(
	rpc blockchain.Rpc,
	contracts blockchain.RskContracts,
	peginQuoteRepository quote.PeginQuoteRepository,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
) *GetQuoteUseCase {
	return &GetQuoteUseCase{
		rpc:                  rpc,
		contracts:            contracts,
		peginQuoteRepository: peginQuoteRepository,
		lp:                   lp,
		peginLp:              peginLp,
	}
}

type QuoteRequest struct {
	callEoaOrContractAddress string
	callContractArguments    []byte
	valueToTransfer          *entities.Wei
	rskRefundAddress         string
}

func NewQuoteRequest(
	callEoaOrContractAddress string,
	callContractArguments []byte,
	valueToTransfer *entities.Wei,
	rskRefundAddress string,
) QuoteRequest {
	return QuoteRequest{
		callEoaOrContractAddress: callEoaOrContractAddress,
		callContractArguments:    callContractArguments,
		valueToTransfer:          valueToTransfer,
		rskRefundAddress:         rskRefundAddress,
	}
}

type GetPeginQuoteResult struct {
	PeginQuote quote.PeginQuote
	Hash       string
}

func (useCase *GetQuoteUseCase) Run(ctx context.Context, request QuoteRequest) (GetPeginQuoteResult, error) {
	var peginQuote quote.PeginQuote
	var creationData quote.PeginCreationData
	var fedAddress string
	var errorArgs usecases.ErrorArgs
	var err error
	var estimatedCallGas *entities.Wei

	if err = usecases.CheckPauseState(useCase.contracts.PegIn); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	peginConfiguration := useCase.peginLp.PeginConfiguration(ctx)
	if errorArgs, err = useCase.validateRequest(peginConfiguration, request); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseErrorArgs(usecases.GetPeginQuoteId, err, errorArgs)
	}

	estimatedCallGas, err = useCase.rpc.Rsk.EstimateGas(ctx, request.callEoaOrContractAddress, request.valueToTransfer, request.callContractArguments)
	if err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if fedAddress, err = useCase.getFederationAddress(); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}

	if creationData, err = useCase.buildCreationData(ctx, peginConfiguration); err != nil {
		return GetPeginQuoteResult{}, err
	}

	generalConfiguration := useCase.lp.GeneralConfiguration(ctx)
	fees := quote.Fees{
		CallFee:    quote.CalculateCallFee(request.valueToTransfer, peginConfiguration),
		GasFee:     new(entities.Wei).Mul(estimatedCallGas, creationData.GasPrice),
		PenaltyFee: peginConfiguration.PenaltyFee,
	}
	if peginQuote, err = useCase.buildPeginQuote(generalConfiguration, peginConfiguration, request, fedAddress, estimatedCallGas, fees); err != nil {
		return GetPeginQuoteResult{}, err
	}

	if err = usecases.ValidateMinLockValue(usecases.GetPeginQuoteId, useCase.contracts.Bridge, peginQuote.Value); err != nil {
		return GetPeginQuoteResult{}, err
	}

	return useCase.storeResult(ctx, peginQuote, creationData)
}

func (useCase *GetQuoteUseCase) storeResult(
	ctx context.Context,
	peginQuote quote.PeginQuote,
	creationData quote.PeginCreationData,
) (GetPeginQuoteResult, error) {
	var hash string
	var err error

	if hash, err = useCase.contracts.PegIn.HashPeginQuote(peginQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	createdQuote := quote.CreatedPeginQuote{Quote: peginQuote, CreationData: creationData, Hash: hash}
	if err = useCase.peginQuoteRepository.InsertQuote(ctx, createdQuote); err != nil {
		return GetPeginQuoteResult{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	return GetPeginQuoteResult{PeginQuote: peginQuote, Hash: hash}, nil
}

func (useCase *GetQuoteUseCase) validateRequest(configuration liquidity_provider.PeginConfiguration, request QuoteRequest) (usecases.ErrorArgs, error) {
	var err error
	args := usecases.NewErrorArgs()
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
	var btcRefundAddress string
	const mainnet = "mainnet"

	if nonce, err = utils.GetRandomInt(); err != nil {
		return quote.PeginQuote{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	if useCase.rpc.Btc.NetworkName() == mainnet {
		btcRefundAddress = blockchain.BitcoinMainnetP2PKHZeroAddress
	} else {
		btcRefundAddress = blockchain.BitcoinTestnetP2PKHZeroAddress
	}

	peginQuote := quote.PeginQuote{
		FedBtcAddress:      fedAddress,
		LbcAddress:         useCase.contracts.PegIn.GetAddress(),
		LpRskAddress:       useCase.lp.RskAddress(),
		BtcRefundAddress:   btcRefundAddress,
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
	}

	if err = entities.ValidateStruct(peginQuote); err != nil {
		return quote.PeginQuote{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	return peginQuote, nil
}

func (useCase *GetQuoteUseCase) getFederationAddress() (string, error) {
	var fedAddress string
	var err error
	if fedAddress, err = useCase.contracts.Bridge.GetFedAddress(); err != nil {
		return "", err
	} else if !blockchain.IsBtcP2SHAddress(fedAddress) {
		return "", errors.New("only P2SH addresses are supported for federation address")
	}
	return fedAddress, nil
}

func (useCase *GetQuoteUseCase) buildCreationData(
	ctx context.Context,
	configuration liquidity_provider.PeginConfiguration,
) (quote.PeginCreationData, error) {
	var gasPrice *entities.Wei
	var err error

	if gasPrice, err = useCase.rpc.Rsk.GasPrice(ctx); err != nil {
		return quote.PeginCreationData{}, usecases.WrapUseCaseError(usecases.GetPeginQuoteId, err)
	}
	creationData := quote.PeginCreationData{
		GasPrice:      gasPrice,
		FeePercentage: configuration.FeePercentage,
		FixedFee:      configuration.FixedFee,
	}
	return creationData, nil
}
