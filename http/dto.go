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
	ApiBaseUrl              string `json:"apiBaseUrl" required:"" example:"https://api.example.com" description:"Provider API base URL"`
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
	FedBTCAddr         string `json:"fedBTCAddr" required:""`
	LBCAddr            string `json:"lbcAddr" required:""`
	LPRSKAddr          string `json:"lpRSKAddr" required:""`
	BTCRefundAddr      string `json:"btcRefundAddr" required:""`
	RSKRefundAddr      string `json:"rskRefundAddr" required:""`
	LPBTCAddr          string `json:"lpBTCAddr" required:""`
	CallFee            uint64 `json:"callFee" required:""`
	PenaltyFee         uint64 `json:"penaltyFee" required:""`
	ContractAddr       string `json:"contractAddr" required:""`
	Data               string `json:"data" required:""`
	GasLimit           uint32 `json:"gasLimit,omitempty" required:""`
	Nonce              int64  `json:"nonce" required:""`
	Value              uint64 `json:"value" required:""`
	AgreementTimestamp uint32 `json:"agreementTimestamp" required:""`
	TimeForDeposit     uint32 `json:"timeForDeposit" required:""`
	LpCallTime         uint32 `json:"lpCallTime" required:""`
	Confirmations      uint16 `json:"confirmations" required:""`
	CallOnRegister     bool   `json:"callOnRegister" required:""`
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
