package pegin

import "errors"

const (
	// MaxPeginDataSize size limit for the data field of the pegin quote
	MaxPeginDataSize = 4_096
	// MaxPeginDepositTxSize size limit allowed for the pegin deposit transaction, the server will reject
	// any pegin deposit transaction that exceeds this limit
	MaxPeginDepositTxSize = 100_000
)

var DataCapExceededError = errors.New("data size exceeds the maximum allowed limit")
