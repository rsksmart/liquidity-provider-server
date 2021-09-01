package models

import "math/big"

type QuoteRequest struct {
	CallContractAddress   string  `json:"callContractAddress"`
	CallContractArguments string  `json:"callContractArguments"`
	ValueToTransfer       big.Int `json:"valueToTransfer"`
	GasLimit              uint    `json:"gasLimit"`
	RskRefundAddress      string  `json:"rskRefundAddress"`
	BitcoinRefundAddress  string  `json:"bitcoinRefundAddress"`
}
