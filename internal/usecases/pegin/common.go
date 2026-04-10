package pegin

import "errors"

const (
	// CallForUserExtraGas
	/**
	 *	The gas spent in the callForUser function can be divided in two parts,
	 *	the first part is the gas spent in the callForUser function itself, call done on behalf
	 *	of the user. This constant represents the first part and needs to be added to the estimation
	 *	done during the get pegin quote process.
	 */
	CallForUserExtraGas = 180000
	// MaxPeginDataSize size limit for the data field of the pegin quote
	MaxPeginDataSize = 4_096
	// MaxPeginDepositTxSize size limit allowed for the pegin deposit transaction, the server will reject
	// any pegin deposit transaction that exceeds this limit
	MaxPeginDepositTxSize = 100_000
)

var DataCapExceededError = errors.New("data size exceeds the maximum allowed limit")
