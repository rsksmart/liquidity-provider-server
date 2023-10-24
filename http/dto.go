package http

import (
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"
)

type ProviderDTO struct {
	Id           uint64 `json:"id" required:"" example:"1" description:"Provider id"`
	Provider     string `json:"provider" required:"" example:"0x0000000000000000000000000000000000000000" description:"Provider address"`
	Name         string `json:"name" required:"" example:"Default Pegin Provider" description:"Provider name"`
	ApiBaseUrl   string `json:"apiBaseUrl" required:"" example:"https://api.example.com" description:"Provider's LPS instance URL"`
	Status       bool   `json:"status" required:"" example:"true" description:"Provider status"`
	ProviderType string `json:"providerType" required:"" example:"pegin" description:"Provider Type"`
}

func toProviderDTO(provider *bindings.LiquidityBridgeContractLiquidityProvider) *ProviderDTO {
	return &ProviderDTO{
		Id:           provider.Id.Uint64(),
		Provider:     provider.Provider.Hex(),
		Name:         provider.Name,
		ApiBaseUrl:   provider.ApiBaseUrl,
		Status:       provider.Status,
		ProviderType: provider.ProviderType,
	}
}
func toGlobalProvider(provider *bindings.LiquidityBridgeContractLiquidityProvider) *types.GlobalProvider {
	return &types.GlobalProvider{
		Id:           provider.Id.Uint64(),
		Provider:     provider.Provider.Hex(),
		Name:         provider.Name,
		ApiBaseUrl:   provider.ApiBaseUrl,
		Status:       provider.Status,
		ProviderType: provider.ProviderType,
	}
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
	CallCost           uint64 `json:"callCost" required:"" description:"The estimated cost for the LP to do the call on behalf of the user. Is calculated with gasPrice * gasLimit"`
}

func toPeginQuote(quote *pegin.Quote) *PeginQuoteDTO {
	return &PeginQuoteDTO{
		FedBTCAddr:         quote.FedBTCAddr,
		LBCAddr:            quote.LBCAddr,
		LPRSKAddr:          quote.LPRSKAddr,
		BTCRefundAddr:      quote.BTCRefundAddr,
		RSKRefundAddr:      quote.RSKRefundAddr,
		LPBTCAddr:          quote.LPBTCAddr,
		CallFee:            quote.CallFee.Uint64(),
		PenaltyFee:         quote.PenaltyFee.Uint64(),
		ContractAddr:       quote.ContractAddr,
		Data:               quote.Data,
		GasLimit:           quote.GasLimit,
		Nonce:              quote.Nonce,
		Value:              quote.Value.Uint64(),
		AgreementTimestamp: quote.AgreementTimestamp,
		TimeForDeposit:     quote.TimeForDeposit,
		LpCallTime:         quote.LpCallTime,
		Confirmations:      quote.Confirmations,
		CallOnRegister:     quote.CallOnRegister,
		CallCost:           quote.CallCost.Uint64(),
	}
}

type PegoutQuoteDTO struct {
	LBCAddr               string `json:"lbcAddress" required:"" validate:"required"`
	LPRSKAddr             string `json:"liquidityProviderRskAddress" required:"" validate:"required"`
	BtcRefundAddr         string `json:"btcRefundAddress" required:"" validate:"required"`
	RSKRefundAddr         string `json:"rskRefundAddress" required:"" validate:"required"`
	LpBTCAddr             string `json:"lpBtcAddr" required:"" validate:"required"`
	CallFee               uint64 `json:"callFee" required:"" validate:"required"`
	PenaltyFee            uint64 `json:"penaltyFee" required:"" validate:"required"`
	Nonce                 int64  `json:"nonce" required:"" validate:"required"`
	DepositAddr           string `json:"depositAddr" required:"" validate:"required"`
	Value                 uint64 `json:"value" required:"" validate:"required"`
	AgreementTimestamp    uint32 `json:"agreementTimestamp" required:"" validate:"required"`
	DepositDateLimit      uint32 `json:"depositDateLimit" required:"" validate:"required"`
	DepositConfirmations  uint16 `json:"depositConfirmations" required:"" validate:"required"`
	TransferConfirmations uint16 `json:"transferConfirmations" required:"" validate:"required"`
	TransferTime          uint32 `json:"transferTime" required:"" validate:"required"`
	ExpireDate            uint32 `json:"expireDate" required:"" validate:"required"`
	ExpireBlock           uint32 `json:"expireBlocks" required:"" validate:"required"`
	CallCost              uint64 `json:"callCost" required:"" description:"The estimated cost for the LP to do the transaction on behalf of the user in Bitcoin network"`
}

func toPegoutQuote(quote *pegout.Quote) *PegoutQuoteDTO {
	return &PegoutQuoteDTO{
		LBCAddr:               quote.LBCAddr,
		LPRSKAddr:             quote.LPRSKAddr,
		BtcRefundAddr:         quote.BtcRefundAddr,
		RSKRefundAddr:         quote.RSKRefundAddr,
		LpBTCAddr:             quote.LpBTCAddr,
		CallFee:               quote.CallFee.Uint64(),
		PenaltyFee:            quote.PenaltyFee,
		Nonce:                 quote.Nonce,
		DepositAddr:           quote.DepositAddr,
		Value:                 quote.Value.Uint64(),
		AgreementTimestamp:    quote.AgreementTimestamp,
		DepositDateLimit:      quote.DepositDateLimit,
		DepositConfirmations:  quote.DepositConfirmations,
		TransferConfirmations: quote.TransferConfirmations,
		TransferTime:          quote.TransferTime,
		ExpireDate:            quote.ExpireDate,
		ExpireBlock:           quote.ExpireBlock,
		CallCost:              quote.CallCost.Uint64(),
	}
}
