package http

import (
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	"github.com/rsksmart/liquidity-provider-server/pegin"
)

type ProviderDTO struct {
	Id                      uint64 `json:"id" required:"" example:"1" description:"Provider id"`
	Provider                string `json:"provider" required:"" example:"0x0000000000000000000000000000000000000000" description:"Provider address"`
	Name                    string `json:"name" required:"" example:"Default Pegin Provider" description:"Provider name"`
	Fee                     uint64 `json:"fee" required:"" example:"100000000000000000" description:"Provider fee"`
	QuoteExpiration         uint64 `json:"quoteExpiration" required:"" example:"3600" description:"Quote expiration time in seconds"`
	AcceptedQuoteExpiration uint64 `json:"acceptedQuoteExpiration" required:"" example:"3600" description:"Accepted quote expiration time in seconds"`
	MinTransactionValue     uint64 `json:"minTransactionValue" required:"" example:"100000000" description:"Minimum transaction value"`
	MaxTransactionValue     uint64 `json:"maxTransactionValue" required:"" example:"1000000000000000000" description:"Maximum transaction value"`
	ApiBaseUrl              string `json:"apiBaseUrl" required:"" example:"https://api.example.com" description:"Provider's LPS instance URL"`
	Status                  bool   `json:"status" required:"" example:"true" description:"Provider status"`
}

func toProviderDTO(provider *bindings.LiquidityBridgeContractProvider) *ProviderDTO {
	return &ProviderDTO{
		Id:                      provider.Id.Uint64(),
		Provider:                provider.Provider.Hex(),
		Name:                    provider.Name,
		Fee:                     provider.Fee.Uint64(),
		QuoteExpiration:         provider.QuoteExpiration.Uint64(),
		AcceptedQuoteExpiration: provider.AcceptedQuoteExpiration.Uint64(),
		MinTransactionValue:     provider.MinTransactionValue.Uint64(),
		MaxTransactionValue:     provider.MaxTransactionValue.Uint64(),
		ApiBaseUrl:              provider.ApiBaseUrl,
		Status:                  provider.Status,
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
	}
}
