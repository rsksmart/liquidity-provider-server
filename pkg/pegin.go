package pkg

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"math/big"
)

type PeginQuoteRequest struct {
	CallEoaOrContractAddress string `json:"callEoaOrContractAddress" required:"" validate:"required,eth_addr" example:"0x0" description:"Contract address or EOA address"`
	CallContractArguments    string `json:"callContractArguments" required:"" validate:"" example:"0x0" description:"Contract data"`
	ValueToTransfer          uint64 `json:"valueToTransfer" required:"" validate:"required" example:"0x0" description:"Value to send in the call"`
	RskRefundAddress         string `json:"rskRefundAddress" required:"" validate:"required,eth_addr" example:"0x0" description:"User RSK refund address"`
}

type PeginQuoteDTO struct {
	FedBTCAddr         string `json:"fedBTCAddr" required:"" description:"The BTC address of the PowPeg"`
	LBCAddr            string `json:"lbcAddr" required:"" description:"The address of the LBC"`
	LPRSKAddr          string `json:"lpRSKAddr" required:"" description:"The RSK address of the LP"`
	BTCRefundAddr      string `json:"btcRefundAddr" required:"" description:"A User BTC refund address"`
	RSKRefundAddr      string `json:"rskRefundAddr" required:"" description:"A User RSK refund address"`
	LPBTCAddr          string `json:"lpBTCAddr" required:"" description:"The BTC address of the LP"`
	CallFee            uint64 `json:"callFee" required:"" description:"The fee charged by the LP"`
	PenaltyFee         uint64 `json:"penaltyFee" required:"" description:"The penalty fee that the LP pays if it fails to deliver the service"`
	ContractAddr       string `json:"contractAddr" required:"" description:"The destination address of the peg-in"`
	Data               string `json:"data" required:"" description:"The arguments to send in the call"`
	GasLimit           uint32 `json:"gasLimit,omitempty" required:"" description:"The gas limit"`
	Nonce              int64  `json:"nonce" required:"" description:"A nonce that uniquely identifies this quote"`
	Value              uint64 `json:"value" required:"" description:"The value to transfer in the call"`
	AgreementTimestamp uint32 `json:"agreementTimestamp" required:"" description:"The timestamp of the agreement"`
	TimeForDeposit     uint32 `json:"timeForDeposit" required:"" description:"The time (in seconds) that the user has to achieve one confirmation on the BTC deposit"`
	LpCallTime         uint32 `json:"lpCallTime" required:"" description:"The time (in seconds) that the LP has to perform the call on behalf of the user after the deposit achieves the number of confirmations"`
	Confirmations      uint16 `json:"confirmations" required:"" description:"The number of confirmations that the LP requires before making the call"`
	CallOnRegister     bool   `json:"callOnRegister" required:"" description:"A boolean value indicating whether the callForUser can be called on registerPegIn"`
	GasFee             uint64 `json:"gasFee" required:"" description:"Fee to pay for the gas of every call done during the pegin (call on behalf of the user and call to the dao fee collector)"`
	ProductFeeAmount   uint64 `json:"productFeeAmount" required:"" description:"The DAO Fee amount"`
}

type RetainedPeginQuoteDTO struct {
	QuoteHash           string   `json:"quoteHash" required:"" description:"32-byte long hash of the quote that acts as a unique identifier"`
	Signature           string   `json:"signature" required:"" description:"Signature of the liquidity provider expressing commitment on the quote"`
	DepositAddress      string   `json:"depositAddress" required:"" description:"BTC derivation address where the user should send the BTC"`
	RequiredLiquidity   *big.Int `json:"requiredLiquidity" required:"" description:"RBTC liquidity that the LP locks to guarantee the service. It is different from the total amount that the user needs to pay."`
	State               string   `json:"state" required:"" description:"Current state of the quote. Possible values are:\n - WaitingForDeposit\n - WaitingForDepositConfirmations\n - TimeForDepositElapsed\n - CallForUserSucceeded\n - CallForUserFailed\n - RegisterPegInSucceeded\n - RegisterPegInFailed"`
	UserBtcTxHash       string   `json:"userBtcTxHash" required:"" description:"The hash of the user's BTC transaction to the derivation address"`
	CallForUserTxHash   string   `json:"callForUserTxHash" required:"" description:"The hash of the RSK transaction to the address requested by the user"`
	RegisterPeginTxHash string   `json:"registerPeginTxHash" required:"" description:"The hash of the RSK transaction where the LP gets his refund and fee"`
}

type PeginCreationDataDTO struct {
	GasPrice      uint64  `json:"gasPrice" required:"" description:"The gas price used to compute the gas fee"`
	FeePercentage float64 `json:"percentageFee" required:"" description:"The percentage fee used to compute the call fee"`
	FixedFee      uint64  `json:"fixedFee" required:"" description:"The fixed fee used to compute the call fee"`
}

type PeginQuoteStatusDTO struct {
	Detail       PeginQuoteDTO         `json:"detail" required:"" description:"Agreed specification of the quote"`
	Status       RetainedPeginQuoteDTO `json:"status" required:"" description:"Current status of the quote"`
	CreationData PeginCreationDataDTO  `json:"creationData" required:"" description:"Values used to compute some fields of the quote"`
}

type GetPeginQuoteResponse struct {
	Quote     PeginQuoteDTO `json:"quote" required:"" description:"Detail of the quote"`
	QuoteHash string        `json:"quoteHash" required:"" description:"This is a 64 digit number that derives from a quote object"`
}

type AcceptPeginRespose struct {
	Signature                 string `json:"signature" required:"" example:"0x0" description:"Signature of the quote"`
	BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash" required:"" example:"0x0" description:"Hash of the deposit BTC address"`
}

func FromPeginQuoteDTO(dto PeginQuoteDTO) quote.PeginQuote {
	return quote.PeginQuote{
		FedBtcAddress:      dto.FedBTCAddr,
		LbcAddress:         dto.LBCAddr,
		LpRskAddress:       dto.LPRSKAddr,
		BtcRefundAddress:   dto.BTCRefundAddr,
		RskRefundAddress:   dto.RSKRefundAddr,
		LpBtcAddress:       dto.LPBTCAddr,
		CallFee:            entities.NewUWei(dto.CallFee),
		PenaltyFee:         entities.NewUWei(dto.PenaltyFee),
		ContractAddress:    dto.ContractAddr,
		Data:               dto.Data,
		GasLimit:           dto.GasLimit,
		Nonce:              dto.Nonce,
		Value:              entities.NewUWei(dto.Value),
		AgreementTimestamp: dto.AgreementTimestamp,
		TimeForDeposit:     dto.TimeForDeposit,
		LpCallTime:         dto.LpCallTime,
		Confirmations:      dto.Confirmations,
		CallOnRegister:     dto.CallOnRegister,
		GasFee:             entities.NewUWei(dto.GasFee),
		ProductFeeAmount:   dto.ProductFeeAmount,
	}
}

func ToPeginQuoteDTO(entity quote.PeginQuote) PeginQuoteDTO {
	return PeginQuoteDTO{
		FedBTCAddr:         entity.FedBtcAddress,
		LBCAddr:            entity.LbcAddress,
		LPRSKAddr:          entity.LpRskAddress,
		BTCRefundAddr:      entity.BtcRefundAddress,
		RSKRefundAddr:      entity.RskRefundAddress,
		LPBTCAddr:          entity.LpBtcAddress,
		CallFee:            entity.CallFee.Uint64(),
		PenaltyFee:         entity.PenaltyFee.Uint64(),
		ContractAddr:       entity.ContractAddress,
		Data:               entity.Data,
		GasLimit:           entity.GasLimit,
		Nonce:              entity.Nonce,
		Value:              entity.Value.Uint64(),
		AgreementTimestamp: entity.AgreementTimestamp,
		TimeForDeposit:     entity.TimeForDeposit,
		LpCallTime:         entity.LpCallTime,
		Confirmations:      entity.Confirmations,
		CallOnRegister:     entity.CallOnRegister,
		GasFee:             entity.GasFee.Uint64(),
		ProductFeeAmount:   entity.ProductFeeAmount,
	}
}

func ToRetainedPeginQuoteDTO(entity quote.RetainedPeginQuote) RetainedPeginQuoteDTO {
	return RetainedPeginQuoteDTO{
		QuoteHash:           entity.QuoteHash,
		Signature:           entity.Signature,
		DepositAddress:      entity.DepositAddress,
		RequiredLiquidity:   entity.RequiredLiquidity.AsBigInt(),
		State:               string(entity.State),
		UserBtcTxHash:       entity.UserBtcTxHash,
		CallForUserTxHash:   entity.CallForUserTxHash,
		RegisterPeginTxHash: entity.RegisterPeginTxHash,
	}
}

func ToPeginCreationDataDTO(entity quote.PeginCreationData) PeginCreationDataDTO {
	feePercentage, _ := entity.FeePercentage.Native().Float64()
	return PeginCreationDataDTO{
		GasPrice:      entity.GasPrice.Uint64(),
		FeePercentage: feePercentage,
		FixedFee:      entity.FixedFee.Uint64(),
	}
}
